#!/bin/bash

docker-compose -f zookeeper2-docker-compose.yaml up -d 
docker-compose -f kafka2-docker-compose.yaml up -d 
docker-compose -f orderer1.docker-compose.yaml up -d 


