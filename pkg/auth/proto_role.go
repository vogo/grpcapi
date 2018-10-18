package auth

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io/ioutil"

	"github.com/golang/protobuf/proto"
	protobuf "github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/vogo/clog"
	"github.com/vogo/grpcapi/api"
	"github.com/vogo/grpcapi/pkg/pb"
)

//MethodRoles roles
type MethodRoles map[pb.Role]bool

//ServiceRoles roles
type ServiceRoles map[string]MethodRoles

var (
	serviceRolesMap = make(map[string]ServiceRoles)
)

//getServiceRoles get
func getServiceRoles(name string) ServiceRoles {
	c, ok := serviceRolesMap[name]
	if !ok {
		c = ServiceRoles{}
		serviceRolesMap[name] = c
	}
	return c
}

//getMethodRoles get
func getMethodRoles(serviceRoles ServiceRoles, name string) MethodRoles {
	c, ok := serviceRoles[name]
	if !ok {
		c = MethodRoles{}
		serviceRoles[name] = c
	}
	return c
}

//AllowRoles whether allow roles
func AllowRoles(serviceName, methodName string, roles []pb.Role) bool {
	serviceRoles, ok := serviceRolesMap[serviceName]
	if !ok {
		return true
	}
	methodRoles, ok := serviceRoles[methodName]

	if !ok {
		return true
	}

	if len(roles) == 0 {
		return false
	}

	for _, role := range roles {
		if _, ok := methodRoles[role]; !ok {
			return false
		}
	}

	return true
}

// copy from: https://github.com/golang/protobuf/blob/master/descriptor/descriptor.go
// extractFile extracts a FileDescriptorProto from a gzip'd buffer.
func extractFile(gz []byte) (*protobuf.FileDescriptorProto, error) {
	r, err := gzip.NewReader(bytes.NewReader(gz))
	if err != nil {
		return nil, fmt.Errorf("failed to open gzip reader: %v", err)
	}
	defer r.Close()

	b, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("failed to uncompress descriptor: %v", err)
	}

	fd := new(protobuf.FileDescriptorProto)
	if err := proto.Unmarshal(b, fd); err != nil {
		return nil, fmt.Errorf("malformed FileDescriptorProto: %v", err)
	}

	return fd, nil
}

//Parse file descriptor for specific filename
func Parse(filename string) *protobuf.FileDescriptorProto {
	b := proto.FileDescriptor(filename)
	if b == nil {
		panic(fmt.Sprintf("invalid proto file name:%s", filename))
	}
	fileDesc, err := extractFile(b)
	if err != nil {
		panic(fmt.Sprintf("invalid FileDescriptorProto: %v", err))
	}
	return fileDesc
}

func getFileAllowRoles(fileDesc *protobuf.FileDescriptorProto) []pb.Role {
	fileOptions := fileDesc.GetOptions()
	roles, err := proto.GetExtension(fileOptions, pb.E_FileAllowRoles)
	if err != nil {
		return nil
	}
	roleArr, ok := roles.([]pb.Role)
	if !ok {
		return nil
	}
	return roleArr
}

func getServiceAllowRoles(service *protobuf.ServiceDescriptorProto) []pb.Role {
	option := service.GetOptions()
	roles, err := proto.GetExtension(option, pb.E_ServiceAllowRoles)
	if err != nil {
		return nil
	}
	roleArr, ok := roles.([]pb.Role)
	if !ok {
		return nil
	}
	return roleArr
}

func getMethodAllowRoles(method *protobuf.MethodDescriptorProto) []pb.Role {
	option := method.GetOptions()
	roles, err := proto.GetExtension(option, pb.E_MethodAllowRoles)
	if err != nil {
		return nil
	}
	roleArr, ok := roles.([]pb.Role)
	if !ok {
		return nil
	}
	return roleArr
}

func parseService(pkg string, fileRoleArr []pb.Role, service *protobuf.ServiceDescriptorProto) {
	name := *service.Name
	serviceName := fmt.Sprintf("%s.%s", pkg, name)
	serviceRoleArr := getServiceAllowRoles(service)

	for _, method := range service.GetMethod() {
		parseMethod(fileRoleArr, serviceRoleArr, serviceName, method)
	}
}
func parseMethod(fileRoleArr []pb.Role, serviceRoleArr []pb.Role, serviceName string, method *protobuf.MethodDescriptorProto) {
	serviceRoles := getServiceRoles(serviceName)
	name := *method.Name
	methodRoles := getMethodRoles(serviceRoles, name)
	for _, role := range fileRoleArr {
		methodRoles[role] = true
	}
	for _, role := range serviceRoleArr {
		methodRoles[role] = true
	}
	for _, role := range getMethodAllowRoles(method) {
		methodRoles[role] = true
	}
}

func parseFile(filename string) {
	fileDesc := Parse(filename)

	fileRoleArr := getFileAllowRoles(fileDesc)

	pkg := fileDesc.GetPackage()
	services := fileDesc.GetService()
	for _, service := range services {
		parseService(pkg, fileRoleArr, service)
	}
}

func init() {
	for _, file := range api.ProtoFiles {
		parseFile(file)
	}
	clog.Info(nil, "service roles map:%v", serviceRolesMap)
}
