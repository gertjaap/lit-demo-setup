#!/bin/bash
set +e
docker stop litdemoadminpanel
docker rm litdemoadminpanel

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

./buildcontainerandrun.sh