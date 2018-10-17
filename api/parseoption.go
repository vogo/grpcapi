package main

import (
	"fmt"

	"github.com/golang/protobuf/descriptor"
	"github.com/vogo/grpcapi/pkg/pb"
)

func main() {
	var msg *pb.TestMessage
	_, md := descriptor.ForMessage(msg)
	options := md.GetOptions()

	fmt.Println(options.String())
	fmt.Println(len(options.GetUninterpretedOption()))
}
