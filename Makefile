start_docker:
	docker compose up -d

start_podman:
	podman-compose up -d

copy_cert:
	cp -r ./cert ./backend/
	cp -r ./cert ./web/