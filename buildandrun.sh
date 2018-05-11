#!/bin/bash
set -e
if [ ! -d "litwebui" ]; then
  ./buildlitwebui.sh
fi

cd admin-ui
npm install
npm run build
rm -rf ../admin-api/static
mv build ../admin-api
mv ../admin-api/build ../admin-api/static 
cd ..
cp -r litwebui/build admin-api/static/
mv admin-api/static/build admin-api/static/litwebui

docker build admin-api -t adminpanel
docker run -v "$PWD/data:/data" -p 8000:8000 -v /var/run/docker.sock:/var/run/docker.sock --name lit-demo-adminpanel --rm adminpanel
