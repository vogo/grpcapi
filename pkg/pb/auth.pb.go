// Code generated by protoc-gen-go. DO NOT EDIT.
// source: auth.proto

package pb

import (
	fmt "fmt"
	math "math"

	proto "github.com/golang/protobuf/proto"
	descriptor "github.com/golang/protobuf/protoc-gen-go/descriptor"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

// role same as that in gitlab
type Role int32

const (
	Role_ADMIN     Role = 0
	Role_OWNER     Role = 1
	Role_MASTER    Role = 2
	Role_DEVELOPER Role = 3
	Role_REPOTER   Role = 4
	Role_USER      Role = 5
)

var Role_name = map[int32]string{
	0: "ADMIN",
	1: "OWNER",
	2: "MASTER",
	3: "DEVELOPER",
	4: "REPOTER",
	5: "USER",
}

var Role_value = map[string]int32{
	"ADMIN":     0,
	"OWNER":     1,
	"MASTER":    2,
	"DEVELOPER": 3,
	"REPOTER":   4,
	"USER":      5,
}

func (x Role) String() string {
	return proto.EnumName(Role_name, int32(x))
}

func (Role) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_8bbd6f3875b0e874, []int{0}
}

var E_FileAllowRoles = &proto.ExtensionDesc{
	ExtendedType:  (*descriptor.FileOptions)(nil),
	ExtensionType: ([]Role)(nil),
	Field:         50000,
	Name:          "grpcapi.file_allow_roles",
	Tag:           "varint,50000,rep,name=file_allow_roles,enum=grpcapi.Role",
	Filename:      "auth.proto",
}

var E_FileAllowScopes = &proto.ExtensionDesc{
	ExtendedType:  (*descriptor.FileOptions)(nil),
	ExtensionType: ([]string)(nil),
	Field:         50001,
	Name:          "grpcapi.file_allow_scopes",
	Tag:           "bytes,50001,rep,name=file_allow_scopes",
	Filename:      "auth.proto",
}

var E_ServiceAllowRoles = &proto.ExtensionDesc{
	ExtendedType:  (*descriptor.ServiceOptions)(nil),
	ExtensionType: ([]Role)(nil),
	Field:         51000,
	Name:          "grpcapi.service_allow_roles",
	Tag:           "varint,51000,rep,name=service_allow_roles,enum=grpcapi.Role",
	Filename:      "auth.proto",
}

var E_ServiceAllowScopes = &proto.ExtensionDesc{
	ExtendedType:  (*descriptor.ServiceOptions)(nil),
	ExtensionType: ([]string)(nil),
	Field:         51001,
	Name:          "grpcapi.service_allow_scopes",
	Tag:           "bytes,51001,rep,name=service_allow_scopes",
	Filename:      "auth.proto",
}

var E_MethodAllowRoles = &proto.ExtensionDesc{
	ExtendedType:  (*descriptor.MethodOptions)(nil),
	ExtensionType: ([]Role)(nil),
	Field:         52000,
	Name:          "grpcapi.method_allow_roles",
	Tag:           "varint,52000,rep,name=method_allow_roles,enum=grpcapi.Role",
	Filename:      "auth.proto",
}

var E_MethodAllowScopes = &proto.ExtensionDesc{
	ExtendedType:  (*descriptor.MethodOptions)(nil),
	ExtensionType: ([]string)(nil),
	Field:         52001,
	Name:          "grpcapi.method_allow_scopes",
	Tag:           "bytes,52001,rep,name=method_allow_scopes",
	Filename:      "auth.proto",
}

func init() {
	proto.RegisterEnum("grpcapi.Role", Role_name, Role_value)
	proto.RegisterExtension(E_FileAllowRoles)
	proto.RegisterExtension(E_FileAllowScopes)
	proto.RegisterExtension(E_ServiceAllowRoles)
	proto.RegisterExtension(E_ServiceAllowScopes)
	proto.RegisterExtension(E_MethodAllowRoles)
	proto.RegisterExtension(E_MethodAllowScopes)
}

func init() { proto.RegisterFile("auth.proto", fileDescriptor_8bbd6f3875b0e874) }

var fileDescriptor_8bbd6f3875b0e874 = []byte{
	// 343 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x84, 0xce, 0xcd, 0x4e, 0xea, 0x40,
	0x14, 0xc0, 0xf1, 0x7b, 0xa1, 0x80, 0x1c, 0x03, 0x0e, 0x83, 0x0b, 0x63, 0x8c, 0xb2, 0x34, 0x2e,
	0x4a, 0xa2, 0x3b, 0x76, 0x18, 0xc6, 0x44, 0x23, 0x94, 0x4c, 0x55, 0x12, 0x13, 0x83, 0x50, 0x06,
	0x68, 0x32, 0x3a, 0x93, 0x4e, 0xd1, 0x37, 0x70, 0xc9, 0x92, 0xb5, 0xbe, 0x85, 0xbe, 0x85, 0xbe,
	0x91, 0x99, 0x4e, 0x31, 0x6d, 0x20, 0x61, 0xd7, 0x8f, 0x73, 0xfe, 0xe7, 0x07, 0x30, 0x98, 0x85,
	0x53, 0x5b, 0x06, 0x22, 0x14, 0xb8, 0x30, 0x09, 0xa4, 0x37, 0x90, 0xfe, 0x7e, 0x6d, 0x22, 0xc4,
	0x84, 0xb3, 0x7a, 0xf4, 0x79, 0x38, 0x1b, 0xd7, 0x47, 0x4c, 0x79, 0x81, 0x2f, 0x43, 0x11, 0x98,
	0xd1, 0x93, 0x0e, 0x58, 0x54, 0x70, 0x86, 0x8b, 0x90, 0x6b, 0xb6, 0xda, 0x97, 0x1d, 0xf4, 0x4f,
	0x3f, 0x3a, 0xbd, 0x0e, 0xa1, 0xe8, 0x3f, 0x06, 0xc8, 0xb7, 0x9b, 0xee, 0x0d, 0xa1, 0x28, 0x83,
	0x4b, 0x50, 0x6c, 0x91, 0x3b, 0x72, 0xed, 0x74, 0x09, 0x45, 0x59, 0xbc, 0x0d, 0x05, 0x4a, 0xba,
	0x8e, 0xfe, 0x67, 0xe1, 0x2d, 0xb0, 0x6e, 0x5d, 0x42, 0x51, 0xae, 0xd1, 0x03, 0x34, 0xf6, 0x39,
	0xeb, 0x0f, 0x38, 0x17, 0xaf, 0xfd, 0x40, 0x70, 0xa6, 0xf0, 0x81, 0x6d, 0x18, 0xf6, 0x92, 0x61,
	0x5f, 0xf8, 0x9c, 0x39, 0x32, 0xf4, 0xc5, 0xb3, 0xda, 0xfb, 0x7e, 0xcb, 0xd6, 0xb2, 0xc7, 0xe5,
	0xd3, 0x92, 0x1d, 0xab, 0x6d, 0x0d, 0xa2, 0x65, 0x9d, 0x69, 0xea, 0x8a, 0x7e, 0x55, 0x8d, 0x2b,
	0xa8, 0x24, 0xc2, 0xca, 0x13, 0x72, 0x63, 0xf9, 0x27, 0x2a, 0x17, 0xe9, 0xce, 0x5f, 0xca, 0x8d,
	0xd6, 0x1a, 0x8f, 0x50, 0x55, 0x2c, 0x78, 0xf1, 0xbd, 0xb4, 0xf3, 0x68, 0xa5, 0xe6, 0x9a, 0xa9,
	0x65, 0xf0, 0x73, 0xbe, 0x96, 0x5a, 0x89, 0x63, 0x09, 0xad, 0x0b, 0xbb, 0xe9, 0x0b, 0x31, 0x78,
	0xe3, 0x89, 0xaf, 0xb9, 0x31, 0xe3, 0x64, 0x33, 0x66, 0x3f, 0x00, 0x7e, 0x62, 0xe1, 0x54, 0x8c,
	0x52, 0xea, 0xc3, 0x95, 0x64, 0x3b, 0x1a, 0x5a, 0x16, 0xdf, 0x17, 0x6b, 0xd1, 0xc8, 0xa4, 0x12,
	0xe6, 0x2e, 0x54, 0x53, 0xf9, 0x98, 0xbc, 0xa9, 0xff, 0xb1, 0x30, 0xe2, 0x4a, 0x22, 0x68, 0xc0,
	0xe7, 0xd6, 0x7d, 0x46, 0x0e, 0x87, 0xf9, 0x68, 0xf1, 0xec, 0x37, 0x00, 0x00, 0xff, 0xff, 0x78,
	0xcc, 0xbc, 0x27, 0xa2, 0x02, 0x00, 0x00,
}
