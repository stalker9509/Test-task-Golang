GOBASE=$(shell cd)
GOBIN=$(GOBASE)
GOFILES=$(wildcard ./cmd/app/*.go)

BINARY_NAME=Test-task-Golang.exe

all: build test

build:
	@go build -o $(GOBIN)\\$(BINARY_NAME) $(GOFILES)

run: build
	@$(GOBIN)\\$(BINARY_NAME)

test:
	@go test -v ./...

docker-build:
	@docker build -t $(BINARY_NAME) .

postgres-up:
	@docker run --name=test-task-db -e POSTGRES_PASSWORD=${DB_PASSWORD} -p 5432:5432 -d postgres

postgres-down:
	@docker stop test-task-db
	@docker rm test-task-db

migrate-up:
	@migrate -path ./migration -database "postgres://postgres:${DB_PASSWORD}@localhost:5432/postgres?sslmode=disable" up

migrate-down:
	@migrate -path ./migration -database "postgres://postgres:${DB_PASSWORD}@localhost:5432/postgres?sslmode=disable" down

clean:
	@go clean
	@del /f $(GOBIN)\\$(BINARY_NAME)

env:
	@copy .envexample .env
