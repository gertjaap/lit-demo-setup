#!/bin/bash
set +e
docker stop litdemoadminpanel
docker rm litdemoadminpanel

set -e

docker build admin-api -t adminpanel
docker run -d -e LITWEBUI=http://localhost:8999/ -p 8000:8000 -v "$PWD/data:/data" -v /var/run/docker.sock:/var/run/docker.sock --name litdemoadminpanel adminpanel
