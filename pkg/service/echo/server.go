// Copyright 2018 The Vogo Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package echo

import (
	"context"
	"fmt"

	"github.com/vogo/grpcapi/pkg/config"
	"github.com/vogo/grpcapi/pkg/pb"
	"github.com/vogo/grpcapi/pkg/server"
	"google.golang.org/grpc"
)

// Server service for echo
type Server struct{}

// Echo echo for request value
func (s *Server) Echo(c context.Context, req *pb.EchoRequest) (res *pb.EchoResponse, err error) {
	res = &pb.EchoResponse{
		Result: Echo(req.Value),
	}

	return res, nil
}

// Serve to start grpc server
func Serve(c *config.Config) {
	server.Serve(fmt.Sprintf(":%d", config.PortEcho), func(s *grpc.Server) {
		pb.RegisterEchoServiceServer(s, &Server{})
	})
}
