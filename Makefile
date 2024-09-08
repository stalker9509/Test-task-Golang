GOBASE=$(shell pwd)
GOBIN=$(GOBASE)
GOFILES=$(wildcard *.go)

BINARY_NAME=Test-task-Golang

all: build test

build:
	@go build -o $(GOBIN)/$(BINARY_NAME) $(GOFILES)

run: build
	@$(GOBIN)/$(BINARY_NAME)

test:
	@go test -v ./...

docker-build:
	@docker build -t $(BINARY_NAME) .

clean:
	@echo "Cleaning..."
	@go clean
	@rm -f $(GOBIN)/$(BINARY_NAME)