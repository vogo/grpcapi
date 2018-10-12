package hello

import (
	"github.com/vogo/grpcapi/pkg/config"
	"github.com/vogo/grpcapi/pkg/conn"
	"github.com/vogo/grpcapi/pkg/pb"
)

//Client hello service grpc client
type Client struct {
	pb.HelloServiceClient
}

//NewClient new hello service grpc client
func NewClient() (*Client, error) {
	c, err := conn.NewClient(config.HostHello, config.PortHello)
	if err != nil {
		return nil, err
	}
	return &Client{
		HelloServiceClient: pb.NewHelloServiceClient(c),
	}, nil
}
