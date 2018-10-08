// Copyright 2018 The Vogo Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package iam

import (
	"context"
	"time"

	"github.com/vogo/grpcapi/pkg/auth"
	"github.com/vogo/grpcapi/pkg/pb"

	"github.com/vogo/grpcapi/pkg/config"
	"github.com/vogo/grpcapi/pkg/server"
	"google.golang.org/grpc"
)

// Server for hello
type Server struct{}

// Token say hello for request value
func (s *Server) Token(c context.Context, req *pb.TokenRequest) (res *pb.TokenResponse, err error) {
	//TODO to validate user password
	expire := 2 * time.Hour
	token, err := auth.Generate(config.SecretKey, expire, req.Username, "admin")
	if err != nil {
		return nil, err
	}

	res = &pb.TokenResponse{
		AccessToken: token,
	}

	return res, nil
}

// Serve to start grpc server
func Serve() {
	server.Serve(config.TokenServiceAddress, func(s *grpc.Server) {
		pb.RegisterTokenManagerServer(s, &Server{})
	})
}
