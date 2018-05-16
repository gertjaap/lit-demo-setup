#!/bin/bash
set -e

cd admin-ui
npm install
npm run build
rm -rf ../admin-api/static
mv build ../admin-api
mv ../admin-api/build ../admin-api/static 
cd ..
docker build admin-api -t adminpanel
docker stop lit-demo-adminpanel
docker rm lit-demo-adminpanel
docker run -e LITWEBUI=http://litwebui.gertjaap.org/ -v "$PWD/data:/data" -v /var/run/docker.sock:/var/run/docker.sock --name lit-demo-adminpanel --network webnetwork adminpanel
