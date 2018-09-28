#!/bin/bash

docker-compose -f kafka3-docker-compose.yaml up -d 
docker-compose -f orderer2.docker-compose.yaml up -d 


