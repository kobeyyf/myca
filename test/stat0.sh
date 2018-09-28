#!/bin/bash

docker-compose -f zookeeper0-docker-compose.yaml up -d

docker-compose -f kafka0-docker-compose.yaml up -d 



