#!/bin/bash

set -e

docker-compose --file config/docker-setup.yml stop
docker-compose --file config/docker-setup.yml rm
docker-compose --file config/docker-setup.yml up -d

while ! mysql -h 127.0.0.1 -P 33306 -u root -pabcdef123456 -e 'SELECT 1' &>/dev/null; do sleep 0.5; done

mysql -h 127.0.0.1 -P 33306 -u root -pabcdef123456 -e 'CREATE DATABASE todoapp;'
