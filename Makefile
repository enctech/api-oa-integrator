start_docker:
	docker compose up -d

start_podman:
	podman-compose up -d

copy_cert:
	cp -r ./cert ./backend/
	cp -r ./cert ./web/

clear_images:
	docker image ls -q --filter "dangling=true" | xargs docker image rm

build_new:
	docker compose up -d --build

run_application:
	make build_new && make clear_images

update_restart:
	git pull
	chmod u+x scripts/db_backup.sh
	make run_application