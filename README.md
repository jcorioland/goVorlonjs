# goVorlonjs
Simple go REST APIs to create [Vorlonjs](https://github.com/microsoftdx/vorlonjs) instances on demand

## How to

### Build the goVorlonjs image

```
docker build -t vorlonjs/govorlonjs -f src/goVorlonjsApi/Dockerfile .
```

### Create a new goVorlonjs container

```
docker run -d -p 8080:8080 -v /var/run/docker.sock:/var/run/docker.sock vorlonjs/govorlonjs:latest
```

### Create a swarm service

```
docker service create goVorlonjs --publish 8080:8080 --mount "type=bind,source=/var/run/docker.sock,target=/var/run/docker.sock" vorlonjs/govorlonjs:latest
```

### Create a new Vorlonjs instances

```
GET http://MACHINE_ENDPOINT:8080/createVorlonContainer?serviceName=NAME_OF_YOUR_SERVICE
```