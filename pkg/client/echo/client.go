package echo

import (
	"github.com/vogo/grpcapi/pkg/config"
	"github.com/vogo/grpcapi/pkg/conn"
	"github.com/vogo/grpcapi/pkg/pb"
)

//Client echo service grpc client
type Client struct {
	pb.EchoServiceClient
}

//NewClient new echo service grpc client
func NewClient() (*Client, error) {
	c, err := conn.NewClient(config.HostEcho, config.PortEcho)
	if err != nil {
		return nil, err
	}
	return &Client{
		EchoServiceClient: pb.NewEchoServiceClient(c),
	}, nil
}
