package main

import (
	"fmt"

	"github.com/golang/protobuf/descriptor"
	"github.com/golang/protobuf/proto"
	"github.com/vogo/grpcapi/pkg/pb"
)

func main() {
	var msg *pb.TestMessage
	_, md := descriptor.ForMessage(msg)
	options := md.GetOptions()

	fmt.Println(options.String())
	fmt.Println(len(options.GetUninterpretedOption()))

	a, _ := proto.GetExtension(options, pb.E_MsgOptionA)
	fmt.Println(*a.(*int32))

	b, _ := proto.GetExtension(options, pb.E_MsgOptionB)
	fmt.Println(*b.(*int32))
}
