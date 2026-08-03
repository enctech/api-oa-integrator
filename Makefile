# ---------------------------------------------------------------------
# Build machine (Mac). Bump the image tags in docker-compose.yaml first,
# then push. Requires: docker login ghcr.io -u <github-user>
# (PAT with write:packages)
# ---------------------------------------------------------------------
push:
	docker buildx bake backend web --push

push_web:
	docker buildx bake web --push

push_backend:
	docker buildx bake backend --push

# Build for linux/amd64 without pushing, to check it compiles
build_images:
	docker buildx bake backend web

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