#!/bin/bash

docker-compose -f zookeeper1-docker-compose.yaml up -d 
docker-compose -f kafka1-docker-compose.yaml up -d 
docker-compose -f orderer0.docker-compose.yaml up -d 


