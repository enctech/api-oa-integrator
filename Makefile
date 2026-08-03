# ---------------------------------------------------------------------
# Publishing happens in GitHub Actions (.github/workflows/build-dwservice.yaml),
# not from here. Packages under ghcr.io/enctech can only be created by
# enctech or by GITHUB_TOKEN acting as this repo; a personal PAT is
# rejected with "denied: permission_denied: create_package".
#
# To publish: bump the tags in docker-compose.yaml, commit, push the
# branch. The workflow builds and pushes both images.
#
# The push targets below are kept for the case where you are authenticated
# as the namespace owner. build_images is the useful one locally - it
# compiles for linux/amd64 without pushing.
# ---------------------------------------------------------------------
BUILDER ?= oa-builder
BAKE = docker buildx --builder $(BUILDER) bake

# Docker Desktop's default builder uses the "docker" driver, which cannot
# export a build cache and cannot build for a foreign platform. Both are
# needed here, so use a docker-container builder instead. Created once,
# then reused; this target is a no-op when it already exists.
builder:
	@docker buildx inspect $(BUILDER) >/dev/null 2>&1 || \
		docker buildx create --name $(BUILDER) --driver docker-container --bootstrap

push: builder
	$(BAKE) backend web --push

push_web: builder
	$(BAKE) web --push

push_backend: builder
	$(BAKE) backend --push

# Build for linux/amd64 without pushing, to check it compiles
build_images: builder
	$(BAKE) backend web

# ---------------------------------------------------------------------
# Server (VM) - pull the prebuilt images and restart.
# ---------------------------------------------------------------------
deploy:
	git pull
	docker compose pull
	docker compose up -d

start_docker:
	docker compose up -d

start_podman:
	podman-compose up -d

copy_cert:
	make copy_cert_backend && make copy_cert_web

copy_cert_web:
	cp -r ./cert ./web/

copy_cert_backend:
	cp -r ./cert ./backend/

clear_images:
	docker image ls -q --filter "dangling=true" | xargs docker image rm

build_new:
	docker compose up -d --build

run_application:
	make update && make build_new && make clear_images

update:
	git pull
	chmod u+x scripts/db_backup.sh
	chmod u+x scripts/startup.sh
	chmod u+x scripts/update_cert.sh

update_restart:
	make update
	make run_application