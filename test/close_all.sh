#!/bin/bash

docker-compose -f peer0-orga-docker-compose.yaml down
docker-compose -f peer0-orgb-docker-compose.yaml down
docker-compose -f peer1-orga-docker-compose.yaml down
docker-compose -f peer1-orgb-docker-compose.yaml down

docker-compose -f orderer0.docker-compose.yaml down
docker-compose -f orderer1.docker-compose.yaml down
docker-compose -f orderer2.docker-compose.yaml down

docker-compose -f kafka0-docker-compose.yaml down
docker-compose -f kafka1-docker-compose.yaml down
docker-compose -f kafka2-docker-compose.yaml down
docker-compose -f kafka3-docker-compose.yaml down

docker-compose -f zookeeper1-docker-compose.yaml down
docker-compose -f zookeeper2-docker-compose.yaml down
docker-compose -f zookeeper0-docker-compose.yaml down

docker-compose -f cli-docker-compose.yaml down
