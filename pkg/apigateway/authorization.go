// Copyright 2018 The Vogo Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package apigateway

import (
	"context"
	"fmt"

	"github.com/vogo/clog"
	"github.com/vogo/grpcapi/pkg/auth"
	"github.com/vogo/grpcapi/pkg/identity"
	"github.com/vogo/grpcapi/pkg/util/ctxutil"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// see https://godoc.org/google.golang.org/grpc#UnaryClientInterceptor
// authorizing in APIGATEWAY only
func authorizer(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	clog.Debug(ctx, "request grpc method:%s", method)

	identityJSON := ctxutil.GetIdentity(ctx)
	identity, err := identity.Parse(identityJSON)
	if err != nil {
		clog.Error(ctx, "failed to parse identity: %v", err)
		return err
	}
	if !auth.AllowRoles(method, identity.Roles) {
		err := status.Error(codes.Unauthenticated, fmt.Sprintf("role %v not allowed to call %s", identity.Roles, method))
		return err
	}
	// Calls the invoker to execute RPC
	err = invoker(ctx, method, req, reply, cc, opts...)
	return err
}
