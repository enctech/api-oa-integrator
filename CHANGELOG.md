# Changelog

## 2026-01-21

### Server Changes

- **Upgraded Grafana Alloy** from v1.6.1 to v1.12.2
  - Fixes PostgreSQL 17 compatibility issue (`column "checkpoints_timed" does not exist`)
  - v1.12.2 has native support for PostgreSQL 17's `pg_stat_checkpointer` view

- **Increased UDP buffer size** for QUIC/HTTP3 performance
  - Added to `/etc/sysctl.conf`:
    ```
    net.core.rmem_max=7500000
    net.core.wmem_max=7500000
    ```
  - Resolves Caddy warning: "failed to sufficiently increase receive buffer size"

### Application Changes

- **Optimized GetLogs query performance** ([PR #2](https://github.com/enctech/api-oa-integrator/pull/2))
  - Added index on `logs.created_at` for faster sorting and range queries
  - Reordered query to filter by date first (uses index)
  - Default date range changed to last 1 hour (prevents full table scan)
  - Configured database connection pool limits

## 2026-01-20

### Infrastructure

- **Upgraded Linode instance** from 1GB to 2GB RAM
