#!/bin/bash

# remove the goVorlonjs API service
docker service rm govorlonjs

# remove the proxy
docker service rm proxy

# remove the swarm listener
docker service rm swarm-listener