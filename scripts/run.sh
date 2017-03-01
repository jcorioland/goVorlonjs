#!/bin/bash

# create vorlonjs network
docker network create -d overlay --attachable vorlonjs

# pull the goVorlonjs API image
docker pull vorlonjs/govorlonjs:0.5.4-dev

# pull the swarm listener
docker pull vfarcic/docker-flow-swarm-listener:1.13

# pull the proxy
docker pull vfarcic/docker-flow-proxy:1.200

# pull the Vorlonjs dashboard image
docker pull vorlonjs/dashboard:0.5.4

# create the swarm listener
docker service create --name swarm-listener \
    --network vorlonjs \
    --mount "type=bind,source=/var/run/docker.sock,target=/var/run/docker.sock" \
    -e DF_NOTIFY_CREATE_SERVICE_URL=http://proxy:8080/v1/docker-flow-proxy/reconfigure \
    -e DF_NOTIFY_REMOVE_SERVICE_URL=http://proxy:8080/v1/docker-flow-proxy/remove \
    --constraint 'node.role==manager' \
    vfarcic/docker-flow-swarm-listener:1.13

# create the proxy
docker service create --name proxy -p 80:80 -p 443:443 \
    --network vorlonjs \
    -e MODE=swarm \
    -e LISTENER_ADDRESS=swarm-listener \
    vfarcic/docker-flow-proxy:1.200

# create the goVorlonjs API service
docker service create --name govorlonjs \
    --network vorlonjs \
    --mount "type=bind,source=/var/run/docker.sock,target=/var/run/docker.sock" \
    --label com.df.notify=true \
    --label com.df.distribute=true \
    --label com.df.servicePath=/api \
    --label com.df.port=82 \
    -e VORLONJS_API_KEY=YOUR_API_KEY \
    -e VORLONJS_API_SECRET=YOUR_API_SECRET \
    vorlonjs/govorlonjs:0.5.4-dev