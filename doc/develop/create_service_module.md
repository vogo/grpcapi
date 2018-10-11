# Create Service Module

## 1. define service api

create a protocol buffer definition file in directory `api/`. eg `api/hello.proto`

then execute `cd api; make` to generate grpc proxy code under directory `pkg/pb/`. eg `pkg/pb/hello.pb.go` and `pkg/pb/hello.pb.gw.go`.

## 2. implement service

create a grpc api implement, eg `pkg/service/hello/server.go` which implements the protocol buffer service interface.

And provid a `Serve()` function to start grpc server
```
// Serve to start grpc server
func Serve(c *config.Config) {
	server.Serve(config.HelloServiceAddress, func(s *grpc.Server) {
		pb.RegisterHelloServiceServer(s, &Server{})
	})
}
```

## 3. service main

create grpc main starter, eg `cmd/hello/main.go`

## 4. Add service proxy into apigateway

add grpc service register map in file `pkg/apigateway/proxy.go`:

```
var (
	registers = []register{
		{config.ServiceKeyEcho, pb.RegisterEchoServiceHandlerFromEndpoint},
		{config.ServiceKeyHello, pb.RegisterHelloServiceHandlerFromEndpoint},
	}
)
```

