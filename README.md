# goVorlonjs
Simple go REST APIs to create [Vorlonjs](https://github.com/microsoftdx/vorlonjs) instances on demand

## How to

### Run the API

#### On Windows (need Docker for Windows)

```
.\scripts\run.ps1
```

#### On Linux (need Docker installed)

```
./run.sh
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
