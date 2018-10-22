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
cd ..

cd pairing-ui
npm install
npm run build
rm -rf ../admin-api/static/pairing
mv build ../admin-api/static
mv ../admin-api/static/build ../admin-api/static/pairing 
cd ..

./buildcontainerandrun.sh