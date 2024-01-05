#!/usr/bin/env bash

rm *.sqlite
docker kill $(docker ps -q)
docker container prune -f
docker network prune -f
docker volume prune -f
docker-compose up -d