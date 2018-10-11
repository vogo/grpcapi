package apigateway

import (
	"context"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/vogo/grpcapi/pkg/config"
	"github.com/vogo/grpcapi/pkg/pb"
	"google.golang.org/grpc"
)

type register struct {
	name string
	f    func(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) (err error)
}

var (
	registers = []register{
		{config.ServiceKeyEcho, pb.RegisterEchoServiceHandlerFromEndpoint},
		{config.ServiceKeyHello, pb.RegisterHelloServiceHandlerFromEndpoint},
	}
)
