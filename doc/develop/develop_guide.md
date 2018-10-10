# Develop Guide

## 1. Preparation

install mongodb server: 
- redhat or centos: https://docs.mongodb.com/manual/tutorial/install-mongodb-on-red-hat/
- install from mongodb server package: https://www.mongodb.com/download-center/v2/community.

create mongodb oauth2 user:
```
use oauth2

db.createUser(
   {
     user: "oauth2",
     pwd: "oauth2",
     roles: [ "readWrite", "dbAdmin" ]
   }
)
```

## 2. Download 

```
go get github.com/vogo/grpcapi

cd $GOPATH/src/github.com/vogo/grpcapi

# enable golang 11 module
export GO111MODULE=on

# download dependency modules
go mod download
```

## 3.Run 

```bash
# Run Echo Grpc Service
go run cmd/echo/main.go

# Run Hello Grpc Service
go run cmd/hello/main.go

# Run ApiGateway Service
go run cmd/apigateway/main.go
```

## 4.Add Mongodb Oauth2 Demo Data

mongodb collection `client` :
```
{ 
    "_id" : "000000", 
    "secret" : "999999", 
    "domain" : "https://localhost", 
    "scopes" : [
        "read", 
        "manage", 
        "admin", 
        "view"
    ], 
    "grant_types" : [
        "password", 
        "refresh_token", 
    ]
}
```

mongodb collection `user` u1 (password:123456):
```
{
    "_id" : "u1",
    "password" : BinData(0, "bTYCH2X8J0a1F9aK1hv//9BYC7/SIyq89OVz3pZvcog="),
    "salt" : BinData(0, "NXl1mO6iV5YpI8YEuy2vTg=="),
    "scopes" : {
        "000000" : "manage,admin"
    }
}
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

