#!/bin/bash
docker build . -t demoadminpanelproxy:latest
docker stop demoadminpanelproxy
docker rm demoadminpanelproxy
docker create --name demoadminpanelproxy --network lit-demo -p 8999:8999 demoadminpanelproxy:latest
docker restart demoadminpanelproxy
