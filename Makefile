# Images for this branch publish with an -internal suffix so they never
# collide with what master's CI pushes.
#
#   make push VERSION=0.13   -> ghcr.io/enctech/{backend,web}-dwservice:0.13-internal
#
# VERSION is deliberately required rather than defaulting to "latest", so a
# deploy always names the exact image it is running.
VERSION ?=
TAG = $(VERSION)-internal
COMPOSE_BUILD = -f docker-compose.yaml -f docker-compose.build.yaml

define need_version
@if [ -z "$(VERSION)" ]; then \
	echo ""; \
	echo "  ERROR: VERSION is required."; \
	echo "     e.g. make $@ VERSION=0.13   ->  :0.13-internal"; \
	echo ""; \
	exit 1; \
fi
endef

# ---------------------------------------------------------------------
# Build machine (Mac) - build for linux/amd64 and push to ghcr.io.
# Requires: docker login ghcr.io -u <github-user>   (PAT, write:packages)
# ---------------------------------------------------------------------
push:
	$(need_version)
	TAG=$(TAG) docker buildx bake $(COMPOSE_BUILD) --push

push_web:
	$(need_version)
	TAG=$(TAG) docker buildx bake $(COMPOSE_BUILD) --push web

push_backend:
	$(need_version)
	TAG=$(TAG) docker buildx bake $(COMPOSE_BUILD) --push backend

# Build without pushing, to check it compiles
build_images:
	$(need_version)
	TAG=$(TAG) docker buildx bake $(COMPOSE_BUILD)

# ---------------------------------------------------------------------
# Server (VM) - pull the prebuilt images and restart. Never builds.
# Records the deployed tag in .env so plain "docker compose ps/logs/down"
# keep working afterwards without re-specifying VERSION.
# ---------------------------------------------------------------------
deploy:
	$(need_version)
	git pull
	echo "TAG=$(TAG)" > .env
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