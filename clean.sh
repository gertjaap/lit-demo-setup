#!/bin/bash
docker kill $(docker ps -a --filter "name=litdemo*" -q)
docker rm $(docker ps -a --filter "name=litdemo*" -q)

sudo rm -rf data/
