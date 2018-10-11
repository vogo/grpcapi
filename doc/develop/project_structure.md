
# Structure

```
├── LICENSE
├── Makefile
├── README.md
├── api      =======> protocol buffer model and api defintion
│   ├── 0.proto
│   ├── Makefile
│   ├── echo.proto
│   └── hello.proto
├── build   ========> build output
│   └── cmd
│       ├── apigateway
│       ├── echo
│       └── hello
├── cmd    =========> command main
│   ├── apigateway
│   │   ├── config.yml ====> config file
│   │   └── main.go    ====> apigateway main
│   ├── echo
│   │   └── main.go
│   └── hello
│       └── main.go
├── doc   ===========> document
│   ├── design ======> design document
│   │   ├── Makefile
│   │   ├── architecture.dot
│   │   ├── architecture.md
│   │   ├── architecture.png
│   │   └── dependency_libraries.md
│   └── develop ======> develop document
│       └── develop_guide.md
├── go.mod
├── go.sum
├── pkg
│   ├── apigateway  ======> apigateway module
│   │   ├── oauth2.go
│   │   ├── server.go
│   │   └── spec
│   │       ├── Makefile
│   │       ├── api.swagger.json
│   │       ├── makestatic.go
│   │       ├── preprocess.jq
│   │       └── static.go
│   ├── auth
│   │   └── auth.go
│   ├── config   =====> configuration module
│   │   ├── config.go
│   │   ├── config_test.go
│   │   └── config_test.yml
│   ├── pb     ========> protocol buffer generated code
│   │   ├── 0.pb.go
│   │   ├── echo.pb.go
│   │   ├── echo.pb.gw.go
│   │   ├── hello.pb.go
│   │   └── hello.pb.gw.go
│   ├── server
│   │   └── server.go
│   └── service
│       ├── echo  ======> echo grpc service module
│       │   ├── server.go
│       │   ├── service.go
│       │   └── service_test.go
│       └── hello ======> hello grpc service module
│           ├── server.go
│           ├── service.go
│           └── service_test.go
├── test
│   ├── client  =====> protocol buffer generated code for client test
│   │   ├── echo_service
│   │   │   ├── echo_parameters.go
│   │   │   ├── echo_responses.go
│   │   │   └── echo_service_client.go
│   │   ├── grpcapi_client.go
│   │   └── hello_service
│   │       ├── hello_parameters.go
│   │       ├── hello_responses.go
│   │       └── hello_service_client.go
│   └── models
│       ├── grpcapi_echo_request.go
│       ├── grpcapi_echo_response.go
│       ├── grpcapi_hello_request.go
│       └── grpcapi_hello_response.go
└── vendor
```

