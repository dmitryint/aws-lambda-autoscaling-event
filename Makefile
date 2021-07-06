-include .env

LAMBDA_HANDLER=aws-lambda-autoscaling-event
#VERSION := $(shell git describe --tags)
VERSION := 0
BUILD := $(shell git rev-parse --short HEAD)
PROJECTNAME := $(shell basename "$(PWD)" | tr -d '[:space:]')

# Go related variables.
GOBASE := $(shell pwd)
#GOPATH := $(GOBASE)/vendor:$(GOBASE)/src
GOBIN := $(GOBASE)/bin
GOFILES := $(wildcard *.go)
GOCACHE := $(GOBASE)/.cache
# Use linker flags to provide version/build settings
LDFLAGS=-ldflags "-X=main.Version=$(VERSION) -X=main.Build=$(BUILD)"

# Make is verbose in Linux. Make it silent.
MAKEFLAGS += --silent

## compile: Compile the binary.
compile:
	@-$(MAKE) -s go-compile 

clean:
	@rm -f bin/*

lambda:
	@echo "  >  Building Lambda packange..."
	@cd bin; zip ../lambda.zip "$(LAMBDA_HANDLER)"

go-compile: go-build

go-build:
	@echo "  >  Building binary..."
	@mkdir -p "$(GOCACHE)"
	@GOOS=darwin GOARCH=amd64 GOCACHE="$(GOCACHE)" go build -mod vendor $(LDFLAGS) -o "$(GOBIN)/$(PROJECTNAME)-darwin-amd64" app/*
	@GOOS=linux GOARCH=amd64 GOCACHE="$(GOCACHE)" go build -mod vendor $(LDFLAGS) -o "$(GOBIN)/$(PROJECTNAME)-linux-amd64" app/*

	@cp "$(GOBIN)/$(PROJECTNAME)-linux-amd64" "$(GOBIN)/$(LAMBDA_HANDLER)"
