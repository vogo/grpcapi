ifeq ($(OS),Windows_NT)
	GOPATH ?= $(USERPROFILE)/go
	GOPATH := /$(subst ;,:/,$(subst \,/,$(subst :,,$(GOPATH))))
	CURDIR := /$(subst :,,$(CURDIR))
	RM := del /q
else
	GOPATH ?= $(HOME)/go
	RM := rm -f
endif

# CMD_FILES=$(sort $(notdir $(wildcard cmd/*)))
SERVICES=apigateway hello echo 
BUILD_ENVS=CGO_ENABLE=0 GOOS=linux
MAIN_DIR=../cmd
CMD_OUTDIR=cmd

# scratch or alpine
DOCKER_FROM=scratch
BUILD_FLAGS=
ifeq ("$(DOCKER_FROM)","scratch")
	BUILD_FLAGS = -a -ldflags '-extldflags "-static"' 
endif

clean: 
	$(RM) $(CMD_OUTDIR)/*

compile: clean
	@echo "============================> compile"
	
	for service in $(SERVICES); \
	do \
		echo "----> compile $$service" ; \
		$(BUILD_ENVS) go build $(BUILD_FLAGS) -o $(CMD_OUTDIR)/$$service $(MAIN_DIR)/$$service/main.go ; \
	done	
	@echo "-----------> compile over"

docker: compile
	@echo "============================> docker build"
	for service in  $(SERVICES); \
	do \
		echo "" ; \
		echo "----> build docker images vogo/grpc-$$service" ; \
		rm -f Dockerfile ; \
		sed "s/<COMMAND_FILE>/$$service/g" Dockerfile.template.$(DOCKER_FROM) > Dockerfile ; \
		docker build -t vogo/grpc-$$service:latest . ;\
	done	
	@echo "-----------> docker build over"

build: docker
	@echo "============================> build over"
