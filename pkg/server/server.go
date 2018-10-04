package server

import (
	"flag"
	"net"

	"github.com/golang/glog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// RegisterServerFunc function to register server
type RegisterServerFunc func(*grpc.Server)

// Serve start a grpc server registering given server
func Serve(address string, register RegisterServerFunc) {
	flag.Parse()
	flag.Set("logtostderr", "true")

	lis, err := net.Listen("tcp", address)
	if err != nil {
		glog.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	register(s)
	//	pb.RegisterEchoServiceServer(s, &echo.Server{})
	// Register reflection service on gRPC server.
	reflection.Register(s)
	glog.Info("start grpc on ", address)
	if err := s.Serve(lis); err != nil {
		glog.Fatalf("failed to serve: %v", err)
	}
}
