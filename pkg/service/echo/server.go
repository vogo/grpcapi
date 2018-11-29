// Copyright 2018 The Vogo Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package echo

import (
	"context"
	"fmt"

	"github.com/vogo/clog"
	"github.com/vogo/grpcapi/pkg/client/hello"
	"github.com/vogo/grpcapi/pkg/config"
	"github.com/vogo/grpcapi/pkg/pb"
	"github.com/vogo/grpcapi/pkg/server"
	"google.golang.org/grpc"
)

// Server service for echo
type Server struct{}

// Echo echo for request value
func (s *Server) Echo(c context.Context, req *pb.EchoRequest) (res *pb.EchoResponse, err error) {
	clog.Info(c, "context: %v", c)
	helloClient, err := hello.NewClient()
	if err != nil {
		return nil, err
	}
	helloReq := &pb.HelloRequest{Name: req.Value}
	helloRes, err := helloClient.Hello(c, helloReq)
	if err != nil {
		return nil, err
	}
	clog.Debug(c, "hello result:%v", helloRes.Result)

	res = &pb.EchoResponse{
		Result: Echo(helloRes.Result),
	}

	return res, nil
}

// Serve to start grpc server
func Serve(c *config.Config) {
	server.Serve(fmt.Sprintf(":%d", config.PortEcho), func(s *grpc.Server) {
		pb.RegisterEchoServiceServer(s, &Server{})
	})
}
