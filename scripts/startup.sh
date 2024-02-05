#!/bin/bash
cd /home/tng/dev/api-oa-integrator || true
# Remove all stopped containers, images, caches...
docker system prune -a -f
docker compose up -d