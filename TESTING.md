# Testing Guide

## Testing the GetLogs Query Optimization

### Prerequisites

- Docker and Docker Compose
- Access to a PostgreSQL client (or use `docker exec`)

### 1. Start Local Environment

```bash
cd /path/to/api-oa-integrator
docker compose up -d
```

### 2. Verify Index Exists

Check if the index was created:

```bash
docker exec db psql -U postgres -c "\d logs"
```

Expected output should include:
```
Indexes:
    "logs_pkey" PRIMARY KEY, btree (id)
    "idx_logs_created_at" btree (created_at DESC)
```

If the index doesn't exist (existing database), create it manually:

```bash
docker exec db psql -U postgres -c "CREATE INDEX CONCURRENTLY idx_logs_created_at ON logs(created_at DESC);"
```

### 3. Test Backend Query Performance

#### Before optimization (without index)

```bash
docker exec db psql -U postgres -c "DROP INDEX IF EXISTS idx_logs_created_at;"
docker exec db psql -U postgres -c "EXPLAIN ANALYZE SELECT * FROM logs WHERE created_at >= NOW() - INTERVAL '1 hour' AND created_at <= NOW() ORDER BY created_at DESC LIMIT 100;"
```

#### After optimization (with index)

```bash
docker exec db psql -U postgres -c "CREATE INDEX idx_logs_created_at ON logs(created_at DESC);"
docker exec db psql -U postgres -c "EXPLAIN ANALYZE SELECT * FROM logs WHERE created_at >= NOW() - INTERVAL '1 hour' AND created_at <= NOW() ORDER BY created_at DESC LIMIT 100;"
```

**Expected improvement:**
- Before: `Seq Scan` or `Parallel Seq Scan` (slow)
- After: `Index Scan Backward` (fast)

### 4. Test API Endpoint

#### Test with default date range (last 1 hour)

```bash
curl -s -w "\nTime: %{time_total}s\n" "http://localhost:1323/api/transactions/logs?page=0&perPage=100"
```

#### Test with explicit date range

```bash
START=$(date -u -v-1H +"%Y-%m-%dT%H:%M:%SZ")  # 1 hour ago (macOS)
END=$(date -u +"%Y-%m-%dT%H:%M:%SZ")          # now

curl -s -w "\nTime: %{time_total}s\n" "http://localhost:1323/api/transactions/logs?startAt=${START}&endAt=${END}&page=0&perPage=100"
```

#### Test with message filter

```bash
curl -s -w "\nTime: %{time_total}s\n" "http://localhost:1323/api/transactions/logs?message=error&page=0&perPage=100"
```

**Expected:** Response time should be under 500ms with the index.

### 5. Test Frontend Default Date Range

1. Start the web app:
   ```bash
   cd web
   npm install
   npm run dev
   ```

2. Open http://localhost:3000/logs in your browser

3. Verify:
   - "Start Date" picker shows time from 1 hour ago
   - "End Date" picker shows current time
   - URL contains `startAt` and `endAt` parameters
   - Logs table loads quickly

4. Clear the date filters and verify:
   - URL still has `startAt` and `endAt` (defaults applied)
   - Page doesn't hang (no full table scan)

### 6. Performance Comparison

Run this on a database with significant data to compare:

```bash
# Without date filter (old behavior - will be slow without index)
time curl -s "http://localhost:1323/api/transactions/logs?page=0&perPage=100" > /dev/null

# With 1 hour filter (new default)
time curl -s "http://localhost:1323/api/transactions/logs?startAt=$(date -u -v-1H +%Y-%m-%dT%H:%M:%SZ)&endAt=$(date -u +%Y-%m-%dT%H:%M:%SZ)&page=0&perPage=100" > /dev/null
```

### Expected Results

| Scenario | Without Index | With Index |
|----------|---------------|------------|
| No date filter (full scan) | 4-6s | <1s |
| 1 hour filter | 1-2s | <100ms |
| 1 hour + message filter | 1-2s | <200ms |
