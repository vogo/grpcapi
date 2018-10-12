// Copyright 2018 The Vogo Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package server

import (
	"context"
	"errors"
	"fmt"
	"net"
	"runtime/debug"
	"strings"
	"time"

	"github.com/golang/glog"
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	"github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/grpc-ecosystem/go-grpc-middleware/validator"
	"github.com/vogo/grpcapi/pkg/util/ctxutil"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

// GRPCServerCallbackFunc function to register server
type GRPCServerCallbackFunc func(*grpc.Server)

// Serve start a grpc server registering given server
func Serve(address string, callback GRPCServerCallbackFunc) {
	lis, err := net.Listen("tcp", address)
	if err != nil {
		glog.Fatalf("failed to listen: %v", err)
	}
	//	s := grpc.NewServer()
	s := grpc.NewServer(grpcOptions()...)
	callback(s)
	// Register reflection service on gRPC server.
	reflection.Register(s)
	glog.Info("start grpc on ", address)
	if err := s.Serve(lis); err != nil {
		glog.Fatalf("failed to serve: %v", err)
	}
}

func grpcOptions() []grpc.ServerOption {
	builtinOptions := []grpc.ServerOption{
		grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
			MinTime:             10 * time.Second,
			PermitWithoutStream: true,
		}),
		grpc_middleware.WithUnaryServerChain(
			grpc_validator.UnaryServerInterceptor(),
			unaryServerLogInterceptor(),
			//func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
			//	if g.checker != nil {
			//		err = g.checker(ctx, req)
			//		if err != nil {
			//			return
			//		}
			//	}

			//	return handler(ctx, req)
			//},
			//func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
			//	if g.builder != nil {
			//		req = g.builder(ctx, req)
			//	}
			//	return handler(ctx, req)
			//},
			grpc_recovery.UnaryServerInterceptor(
				grpc_recovery.WithRecoveryHandler(func(p interface{}) error {
					glog.Errorf("GRPC server recovery with error: %+v", p)
					glog.Errorf(string(debug.Stack()))
					if e, ok := p.(error); ok {
						return e
					}
					return errors.New("internal error")

				}),
			),
		),
		grpc_middleware.WithStreamServerChain(
			grpc_recovery.StreamServerInterceptor(
				grpc_recovery.WithRecoveryHandler(func(p interface{}) error {
					glog.Errorf("GRPC server recovery with error: %+v", p)
					glog.Errorf(string(debug.Stack()))
					if e, ok := p.(error); ok {
						return e
					}
					return errors.New("internal error")
				}),
			),
		),
	}

	return builtinOptions
}

var (
	jsonPbMarshaller = &jsonpb.Marshaler{
		OrigName: true,
	}
)

func unaryServerLogInterceptor() grpc.UnaryServerInterceptor {

	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		var err error
		s := ctxutil.GetUserID(ctx)
		requestID := ctxutil.GetRequestID(ctx)
		ctx = ctxutil.SetRequestID(ctx, requestID)
		ctx = ctxutil.SetUserID(ctx, s)

		method := strings.Split(info.FullMethod, "/")
		action := method[len(method)-1]

		logPrefix := fmt.Sprintf("%s | %s | %+v |", requestID, action, s)

		if p, ok := req.(proto.Message); ok {
			if content, err := jsonPbMarshaller.MarshalToString(p); err != nil {
				glog.Errorf("%s failed to marshal proto message to string [%+v]", logPrefix, err)
			} else {
				glog.Infof("%s request received: %s", logPrefix, content)
			}
		}
		start := time.Now()

		resp, err := handler(ctx, req)

		elapsed := time.Since(start)
		glog.Infof("%s request elapse: %s", logPrefix, elapsed)
		if e, ok := status.FromError(err); ok {
			if e.Code() != codes.OK {
				glog.V(1).Infof("%s response error: %s, %s", logPrefix, e.Code().String(), e.Message())
			}
		}
		return resp, err
	}
}
