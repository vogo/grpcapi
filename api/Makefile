GOPATH:=$(shell go env GOPATH)
PWD:=$(shell pwd)

PROTO_FILES=$(sort $(wildcard ./*.proto))
PROTOC_INC_PATH=/usr/local/include
GRPC_GATEWAY_PATH=$(GOPATH)/src/github.com/grpc-ecosystem/grpc-gateway
GOOGLEAPIS_PATH=$(GRPC_GATEWAY_PATH)/third_party/googleapis
VENDOR_PATH=../vendor/

PROTOC_FLAGS:=-I. -I$(GOOGLEAPIS_PATH) -I$(GRPC_GATEWAY_PATH) -I$(PROTOC_INC_PATH) -I$(VENDOR_PATH)


generate: $(PROTO_FILES) Makefile
	mkdir -p ../pkg/pb
	rm -rf ../pkg/pb/*

	protoc $(PROTOC_FLAGS) --go_out=plugins=grpc:../pkg/pb ${PROTO_FILES}
	protoc $(PROTOC_FLAGS) --grpc-gateway_out=logtostderr=true,allow_delete_body=true:../pkg/pb ${PROTO_FILES}
	protoc $(PROTOC_FLAGS) --swagger_out=logtostderr=true,allow_delete_body=true:../pkg/apigateway/spec ${PROTO_FILES}

	cd ../pkg/apigateway/spec/ && cat ./*.swagger.json 0.swagger.json | jq --slurp 'reduce .[] as $$item ({}; . * $$item)' | jq -f ./preprocess.jq > ./api.swagger.json
	cd ../pkg/apigateway/spec/ && find . | grep .swagger.json | grep -v "api" | xargs rm

	rm -rf ../test/models/* ../test/client/*
	cd ../test/ && swagger generate client -f ../pkg/apigateway/spec/api.swagger.json -A grpcapi
	# need to exec chown finally
	chown -R $${USER_ID}:$${GROUP_ID} ../pkg ../test
	# format generated files
	goimports -v -l -w -e -local=grpcapi ../pkg/pb
	@echo "ok"

clean:
	@echo "ok"


