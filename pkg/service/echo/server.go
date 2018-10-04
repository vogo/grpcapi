package echo

import (
	"context"

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
func Serve() {
	server.Serve(config.EchoServiceAddress, func(s *grpc.Server) {
		pb.RegisterEchoServiceServer(s, &Server{})
	})
}
