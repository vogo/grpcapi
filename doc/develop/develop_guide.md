# Develop Guide

## 1. Download 

```
go get github.com/vogo/grpcapi

cd $GOPATH/src/github.com/vogo/grpcapi

# enable golang 11 module
export GO111MODULE=on

# download dependency modules
go mod download
```

## 2. Start Mongodb

see [create mongodb](create_mongodb.md)

## 3.Run in Localhost 

edit `/etc/hosts` to add ip-host map for localhost:
```
127.0.0.1 grpc-echo grpc-hello grpc-mongodb
```
> NOTE: clean http proxy setting for shell before starting services.

run services and apigateway:
```bash
# Run Echo Grpc Service
go run cmd/echo/main.go

# Run Hello Grpc Service
go run cmd/hello/main.go

# Run ApiGateway Service
go run cmd/apigateway/main.go -config=cmd/apigateway/config.yml
```
## 4. Run in Docker

create a docker network: `docker network create grpc`

build docker images:
```
cd build
make build
```

start docker instances:
```
cd script
./run-docker-apigateway.sh
./run-docker-echo.sh
./run-docker-hello.sh
```

## 5. Test

request token:
```
curl -X POST \
  http://localhost:8080/oauth2/token \
  -H 'Authorization: Basic MDAwMDAwOjk5OTk5OQ==' \
  -H 'Cache-Control: no-cache' \
  -H 'Content-Type: application/x-www-form-urlencoded' \
  -d 'grant_type=password&username=u1&password=123456&scope=manage%2Cadmin'

# response: {"access_token":"GW8MIJVJN_SNU4DQNYVZIG","expires_in":7200,"refresh_token":"1PO9OUDIU9MLSTCXPNVGKW","scope":"manage,admin","token_type":"Bearer"}

```

request hello service:
```
curl -X POST \
  http://localhost:8080/api/v1/hello \
  -H 'Authorization: Bearer RODJIKFIODABXDJ8I8VLRA' \
  -H 'Cache-Control: no-cache' \
  -H 'Content-Type: application/json' \
  -d '{
	"name": "wongoo"
}'

# response: {"result":"Hello wongoo"}
```

