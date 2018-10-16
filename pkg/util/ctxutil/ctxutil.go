package ctxutil

import (
	"context"

	"github.com/vogo/grpcapi/pkg/auth"
	"google.golang.org/grpc/metadata"
)

type getMetadataFromContext func(ctx context.Context) (md metadata.MD, ok bool)

var getMetadataFromContextFunc = []getMetadataFromContext{
	metadata.FromIncomingContext,
	metadata.FromOutgoingContext,
}

//GetValueFromContext get value from incomming/outging metadata or value of context
func GetValueFromContext(ctx context.Context, key string) []string {
	if ctx == nil {
		return []string{}
	}
	for _, f := range getMetadataFromContextFunc {
		md, ok := f(ctx)
		if !ok {
			continue
		}
		m, ok := md[key]
		if ok && len(m) > 0 {
			return m
		}
	}
	m, ok := ctx.Value(key).([]string)
	if ok && len(m) > 0 {
		return m
	}
	s, ok := ctx.Value(key).(string)
	if ok && len(s) > 0 {
		return []string{s}
	}
	return []string{}
}

//GetRequestID get request id
func GetRequestID(ctx context.Context) string {
	return GetSingleValue(ctx, auth.KeyRequestID)
}

//GetUserID get user id
func GetUserID(ctx context.Context) string {
	return GetSingleValue(ctx, auth.KeyUserID)
}

//GetSingleValue get single value from context
func GetSingleValue(ctx context.Context, key string) string {
	v := GetValueFromContext(ctx, key)
	if len(v) == 0 {
		return ""
	}
	return v[0]
}

//SetUserID set user id
func SetUserID(ctx context.Context, userID string) context.Context {
	ctx = WithUserID(ctx, userID)
	return SetOutgoingContext(ctx, auth.KeyUserID, userID)
}

//WithUserID set user id
func WithUserID(ctx context.Context, userID string) context.Context {
	return WithValue(ctx, auth.KeyUserID, userID)
}

//SetRequestID set request id
func SetRequestID(ctx context.Context, requestID string) context.Context {
	ctx = WithRequestID(ctx, requestID)
	return SetOutgoingContext(ctx, auth.KeyRequestID, requestID)
}

//WithRequestID with request id
func WithRequestID(ctx context.Context, requestID string) context.Context {
	return WithValue(ctx, auth.KeyRequestID, requestID)
}

//WithValue set key/value map into context
func WithValue(ctx context.Context, k interface{}, v string) context.Context {
	return context.WithValue(ctx, k, []string{v})
}

//SetValue set key/value map into context
func SetValue(ctx context.Context, value map[interface{}]string) context.Context {
	for k, v := range value {
		ctx = context.WithValue(ctx, k, []string{v})
		if s, ok := k.(string); ok {
			ctx = SetOutgoingContext(ctx, s, v)
		}
	}
	return ctx
}

//SetOutgoingContext set outgoing context value
func SetOutgoingContext(ctx context.Context, key, value string) context.Context {
	md, ok := metadata.FromOutgoingContext(ctx)
	if !ok {
		md = metadata.MD{}
	}
	md[key] = []string{value}
	return metadata.NewOutgoingContext(ctx, md)
}
