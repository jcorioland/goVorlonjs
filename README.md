# goVorlonjs
Simple go REST APIs to create [Vorlonjs](https://github.com/microsoftdx/vorlonjs) instances on demand

## How to

### Run the API

#### On Windows

```
.\scripts\run.ps1
```

#### On Linux

```
./scripts/run.sh
```

#### Use Docker Compose

```
docker stack deploy govorlonjs -c docker-compose.yml
```

### Use the API

#### Create a new Vorlonjs instance

*Request*

```
POST /api/instance/create HTTP/1.1
Host: REPLACE_WITH_YOUR_HOST
Content-Type: application/json

{
    "serviceName":"SERVICE_NAME"
}
```

*Response*

```
HTTP 201 Created
Vorlonjs is running at /SERVICE_NAME
```

Then your instance is available on `/SERVICE_NAME`

#### Remove an existing Vorlonjs instance

*Request*

```
POST /api/instance/remove HTTP/1.1
Host: REPLACE_WITH_YOUR_HOST
Content-Type: application/json

{
    "serviceName":"vorlonjs1"
}
```

*Response*

```
HTTP 200 Ok
```

#### Check if a Vorlonjs instance exists

*Request*

```
POST /api/instance/exists HTTP/1.1
Host: REPLACE_WITH_YOUR_HOST
Content-Type: application/json

{
    "serviceName":"vorlonjs1"
}
```

*Response*

If the instance exists

```
HTTP 200 Ok
```

If the instance does not exists

```
HTTP 404 Not Found
```

### Clean everything

#### On Windows

```
.\scripts\clean.ps1
```

#### On Linux

```
./scripts/clean.sh
```