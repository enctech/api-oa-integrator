start_docker:
	docker compose up -d

start_podman:
	podman-compose up -d

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

update_restart:
	make update
	make run_application