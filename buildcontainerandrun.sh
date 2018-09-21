#!/bin/bash
docker build admin-api -t adminpanel
docker run -v "$PWD/data:/data" -p 8000:8000 -v /var/run/docker.sock:/var/run/docker.sock --name lit-demo-adminpanel --rm adminpanel
