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

- **Temporarily disabled logs query** ([PR #3](https://github.com/enctech/api-oa-integrator/pull/3))
  - Disabled logs query endpoint to prevent server overload
  - Disabled logging to database
  - Shows maintenance message on logs page

## 2026-01-20

### Infrastructure

- **Upgraded Linode instance** from 1GB to 2GB RAM
