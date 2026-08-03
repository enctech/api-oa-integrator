# CLAUDE.md

## Build & Deploy (chore/TNGdwservice)

Images are built on the Mac and pulled by the VM, which is too slow to
build `react-scripts` itself. Services are `backend`, `web`, `db`, `nginx`.

On the Mac — bump the tags in `docker-compose.yaml`, then:
```bash
docker login ghcr.io -u <github-user>   # PAT with write:packages, once
make push                               # docker buildx bake backend web --push
```

On the VM:
```bash
make deploy                             # git pull && docker compose pull && up -d
```

Builds are pinned to `linux/amd64` in `docker-compose.yaml`. A MacBook
builds arm64 by default, which fails on the VM with `exec format error`.

## Version Tags

Update versions in `docker-compose.yaml` before building:
- `backend-dwservice:X.XX-internal`
- `web-dwservice:X.XX-internal`

The `-internal` suffix keeps these distinct from the images master's CI
publishes as `api-oa-integrator-{backend,web}`.

## Deployment notes

This branch terminates TLS at nginx (port 3030, 80 redirects to it) using
certs mounted from `./ssl`, which is gitignored. Create it before the
first `up`:
```bash
mkdir -p ssl && cp web/cert/certificate.pem web/cert/private-key.pem ssl/
```

`db` is pinned to `postgres:16` — the data directory at
`/var/lib/postgresql/data` was created by 16 and a newer major will not
start on it without `pg_upgrade`.

Caddy rarely changes - only rebuild when Caddyfile is modified.
