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

	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	"github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/grpc-ecosystem/go-grpc-middleware/validator"
	"github.com/vogo/clog"
	"github.com/vogo/grpcapi/pkg/identity"
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
		clog.Fatal(nil, "failed to listen: %v", err)
	}
	//	s := grpc.NewServer()
	s := grpc.NewServer(grpcOptions()...)
	callback(s)
	// Register reflection service on gRPC server.
	reflection.Register(s)
	clog.Info(nil, "start grpc on %s", address)
	if err := s.Serve(lis); err != nil {
		clog.Fatal(nil, "failed to serve: %v", err)
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
					clog.Error(nil, "GRPC server recovery with error: %+v", p)
					clog.Error(nil, string(debug.Stack()))
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
					clog.Error(nil, "GRPC server recovery with error: %+v", p)
					clog.Error(nil, string(debug.Stack()))
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
		start := time.Now()

		var err error
		identityJSON := ctxutil.GetIdentity(ctx)
		requestID := ctxutil.GetRequestID(ctx)
		identity, err := identity.Parse(identityJSON)
		if err != nil {
			clog.Error(ctx, "failed to parse identity: %v", err)
			return nil, err
		}
		method := strings.Split(info.FullMethod, "/")
		action := method[len(method)-1]

		logPrefix := fmt.Sprintf("%s %s %s | %+v |", action, identity.Roles, identity.Scopes, identity.UserID)

		if clog.DebugEnabled() {
			if p, ok := req.(proto.Message); ok {
				if content, err := jsonPbMarshaller.MarshalToString(p); err != nil {
					clog.Error(ctx, "%s failed to marshal proto message to string [%+v]", logPrefix, err)
				} else {
					clog.Debug(ctx, "%s request received: %s", logPrefix, content)
				}
			}
		}

		ctx = ctxutil.SetRequestID(ctx, requestID)
		ctx = ctxutil.SetIdentity(ctx, identityJSON)

		resp, err := handler(ctx, req)
		elapsed := time.Since(start)
		clog.Info(ctx, "%s request elapse: %s", logPrefix, elapsed)

		if e, ok := status.FromError(err); ok {
			if e.Code() != codes.OK {
				clog.Info(ctx, "%s response error: %s, %s", logPrefix, e.Code().String(), e.Message())
			}
		}
		return resp, err
	}
}
