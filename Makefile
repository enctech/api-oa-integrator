TAG ?= latest
COMPOSE_BUILD = -f docker-compose.yaml -f docker-compose.build.yaml

# ---------------------------------------------------------------------
# Build machine (Mac) - build for linux/amd64 and push to ghcr.io.
# Requires: docker login ghcr.io -u <github-user>   (PAT, write:packages)
# ---------------------------------------------------------------------
push:
	TAG=$(TAG) docker buildx bake $(COMPOSE_BUILD) --push

push_web:
	TAG=$(TAG) docker buildx bake $(COMPOSE_BUILD) --push web

push_backend:
	TAG=$(TAG) docker buildx bake $(COMPOSE_BUILD) --push backend

# Build without pushing, to check it compiles
build_images:
	TAG=$(TAG) docker buildx bake $(COMPOSE_BUILD)

# ---------------------------------------------------------------------
# Server (VM) - pull the prebuilt images and restart. Never builds.
# ---------------------------------------------------------------------
deploy:
	git pull
	TAG=$(TAG) docker compose pull
	TAG=$(TAG) docker compose up -d

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