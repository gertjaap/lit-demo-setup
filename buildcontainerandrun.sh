#!/bin/bash
set +e
docker stop litdemoadminpanel
docker rm litdemoadminpanel

set -e

docker build admin-api -t adminpanel
./run.sh