package conn

import (
	"context"
	"crypto/tls"
	"fmt"
	"sync"
	"time"

	"github.com/golang/glog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/keepalive"
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

var clientCache sync.Map

func NewClient(host string, port int) (*grpc.ClientConn, error) {
	endpoint := fmt.Sprintf("%s:%d", host, port)

	if conn, ok := clientCache.Load(endpoint); ok {
		return conn.(*grpc.ClientConn), nil
	}
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, endpoint, ClientOptions...)
	if err != nil {
		glog.V(1).Infof("dial error:%v", err)
		return nil, err
	}
	clientCache.Store(endpoint, conn)
	return conn, nil
}

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
