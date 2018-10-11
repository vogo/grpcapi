// Copyright 2018 The Vogo Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package apigateway

import (
	"context"
	"fmt"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/vogo/grpcapi/pkg/config"
	"github.com/vogo/grpcapi/pkg/pb"
	"google.golang.org/grpc"
)

type register struct {
	endpoint string
	f        func(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) (err error)
}

var (
	registers = []register{
		{fmt.Sprintf("%s%s:%d", config.HostPrefix, config.ServiceEcho, config.PortEcho), pb.RegisterEchoServiceHandlerFromEndpoint},
		{fmt.Sprintf("%s%s:%d", config.HostPrefix, config.ServiceHello, config.PortHello), pb.RegisterHelloServiceHandlerFromEndpoint},
	}
)
