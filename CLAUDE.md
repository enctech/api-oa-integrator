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
- `ghcr.io/yuzurihaaa/backend-dwservice:X.XX-internal`
- `ghcr.io/yuzurihaaa/web-dwservice:X.XX-internal`

The `-internal` suffix keeps these distinct from the images master's CI
publishes as `ghcr.io/enctech/api-oa-integrator-{backend,web}`.

These live under a personal namespace rather than `enctech` because
creating a new package in an org namespace is a separate permission from
`write:packages`, and a personal PAT does not have it. Master's images
exist under `enctech` only because Actions created them with
`GITHUB_TOKEN`. To move these back to `enctech` later, an org owner needs
to allow package creation (and authorise the PAT for SSO if enforced).

The packages are private on first push. The VM therefore needs its own
`docker login ghcr.io` with a PAT carrying `read:packages`, or the
packages can be made public from their GitHub page.

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
