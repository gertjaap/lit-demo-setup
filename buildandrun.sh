#!/bin/bash
set +e
docker stop lit-demo-adminpanel
docker rm lit-demo-adminpanel

set -e

cd admin-ui
npm install
npm run build
rm -rf ../admin-api/static/admin
mv build ../admin-api/static
mv ../admin-api/static/build ../admin-api/static/admin 
cd ../demo-ui
npm install
npm run build
rm -rf ../admin-api/static/demo
mv build ../admin-api/static
mv ../admin-api/static/build ../admin-api/static/demo 
cd ..
docker build admin-api -t adminpanel

docker run -d -e LITWEBUI=http://localhost:8999/ -p 8000:8000 -v "$PWD/data:/data" -v /var/run/docker.sock:/var/run/docker.sock --name lit-demo-adminpanel adminpanel
