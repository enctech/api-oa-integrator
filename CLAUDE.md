# CLAUDE.md

## Build & Deploy

Build only backend and web (skip caddy):
```bash
ssh enctech "cd /root/api-oa-integrator && docker compose build backend web && docker compose up -d"
```

Or with docker bake (faster, parallel builds):
```bash
ssh enctech "cd /root/api-oa-integrator && docker buildx bake backend web && docker compose up -d"
```

## Version Tags

Update versions in `docker-compose.yaml` before building:
- `backend:X.XX`
- `web:X.XX`

Caddy rarely changes - only rebuild when Caddyfile is modified.
