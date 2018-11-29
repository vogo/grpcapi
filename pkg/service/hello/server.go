// Copyright 2018 The Vogo Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package hello

import (
	"context"
	"fmt"

	"github.com/vogo/clog"
	"github.com/vogo/grpcapi/pkg/config"
	"github.com/vogo/grpcapi/pkg/pb"
	"github.com/vogo/grpcapi/pkg/server"
	"google.golang.org/grpc"
)

// Server for hello
type Server struct{}

// Hello say hello for request value
func (s *Server) Hello(c context.Context, req *pb.HelloRequest) (res *pb.HelloResponse, err error) {
	clog.Info(c, "context: %v", c)
	res = &pb.HelloResponse{
		Result: Hello(req.Name),
	}

	return res, nil
}

// Serve to start grpc server
func Serve(c *config.Config) {
	server.Serve(fmt.Sprintf(":%d", config.PortHello), func(s *grpc.Server) {
		pb.RegisterHelloServiceServer(s, &Server{})
	})
}
