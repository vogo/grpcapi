package client

import (
	"context"
	"crypto/tls"
	"fmt"
	"sync"
	"time"

	"github.com/vogo/clog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/keepalive"

	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	"github.com/opentracing/opentracing-go"
)

//ClientOptions grpc client options
var ClientOptions = []grpc.DialOption{
	grpc.WithInsecure(),
	grpc.WithKeepaliveParams(keepalive.ClientParameters{
		Time:                30 * time.Second,
		Timeout:             10 * time.Second,
		PermitWithoutStream: true,
	}),
}

//InitGrpcClientOptions init
func InitGrpcClientOptions() {
	tracingIntercepter := grpc.WithUnaryInterceptor(otgrpc.OpenTracingClientInterceptor(opentracing.GlobalTracer()))
	streamIntercepter := grpc.WithStreamInterceptor(otgrpc.OpenTracingStreamClientInterceptor(opentracing.GlobalTracer()))
	ClientOptions = append(ClientOptions, tracingIntercepter, streamIntercepter)
}

var clientCache sync.Map

//NewClient new grpc client
func NewClient(host string, port int) (*grpc.ClientConn, error) {
	endpoint := fmt.Sprintf("%s:%d", host, port)

	if conn, ok := clientCache.Load(endpoint); ok {
		clientConn := conn.(*grpc.ClientConn)
		if clientConn.GetState() != connectivity.TransientFailure && clientConn.GetState() != connectivity.Shutdown {
			return clientConn, nil
		}
	}
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, endpoint, ClientOptions...)
	if err != nil {
		clog.Debug(nil, "dial error:%v", err)
		return nil, err
	}
	clientCache.Store(endpoint, conn)
	return conn, nil
}

//NewTLSClient new grpc tls cliet
func NewTLSClient(host string, port int, tlsConfig *tls.Config) (*grpc.ClientConn, error) {
	endpoint := fmt.Sprintf("%s:%d", host, port)
	if conn, ok := clientCache.Load(endpoint); ok {
		return conn.(*grpc.ClientConn), nil
	}
	creds := credentials.NewTLS(tlsConfig)
	tlsClientOptions := []grpc.DialOption{
		grpc.WithTransportCredentials(creds),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                30 * time.Second,
			Timeout:             10 * time.Second,
			PermitWithoutStream: true,
		}),
	}
	conn, err := grpc.Dial(endpoint, tlsClientOptions...)
	if err != nil {
		return nil, err
	}
	clientCache.Store(endpoint, conn)
	return conn, nil
}
