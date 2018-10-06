// Copyright 2018 The Vogo Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package apigateway

import (
	"context"
	"flag"
	"net/http"

	"github.com/golang/glog"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/vogo/grpcapi/pkg/apigateway/spec"
	"github.com/vogo/grpcapi/pkg/config"
	"github.com/vogo/grpcapi/pkg/pb"
	"google.golang.org/grpc"
)

type register struct {
	f        func(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) (err error)
	endpoint string
}

var (
	registers = []register{
		{pb.RegisterEchoServiceHandlerFromEndpoint, config.EchoServiceAddress},
		{pb.RegisterHelloServiceHandlerFromEndpoint, config.HelloServiceAddress},
	}

	patternSwagger = runtime.MustPattern(runtime.NewPattern(1, []int{2, 0, 2, 1, 2, 2}, []string{"swagger", "api", "v1"}, ""))
)

func run(address string) error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}

	for _, r := range registers {
		err := r.f(ctx, mux, r.endpoint, opts)
		if err != nil {
			return err
		}
	}

	mux.Handle("GET", patternSwagger, func(w http.ResponseWriter, req *http.Request, pathParams map[string]string) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(spec.Files["api.swagger.json"]))
	})

	return http.ListenAndServe(address, mux)
}

// Serve to start api gateway
func Serve() {
	flag.Parse()
	defer glog.Flush()

	if err := run(config.APIGatewayAddress); err != nil {
		glog.Fatal(err)
	}
}
