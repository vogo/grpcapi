ifeq ($(OS),Windows_NT)
	GOPATH ?= $(USERPROFILE)/go
	GOPATH := /$(subst ;,:/,$(subst \,/,$(subst :,,$(GOPATH))))
	CURDIR := /$(subst :,,$(CURDIR))
	RM := del /q
else
	GOPATH ?= $(HOME)/go
	RM := rm -f
endif

MAKEFILE = $(word $(words $(MAKEFILE_LIST)),$(MAKEFILE_LIST))
BUILD_PARAM = GOOS=linux

clean: 
	$(RM)  build/cmd/*
	mkdir -p build/cmd 

build: clean
	$(BUILD_PARAM) go build -o build/cmd/apigateway cmd/apigateway/main.go 
	$(BUILD_PARAM) go build -o build/cmd/hello cmd/hello/main.go 
	$(BUILD_PARAM) go build -o build/cmd/echo cmd/echo/main.go 


