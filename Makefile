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
	chmod u+x scripts/db_backup.sh
	make build_new && make clear_images

update:
	git pull
	chmod u+x scripts/db_backup.sh
	chmod u+x scripts/startup.sh

update_restart:
	make update
	make run_application