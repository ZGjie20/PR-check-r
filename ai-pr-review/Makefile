.PHONY: run-api run-cli test build

run-api:
	go run ./cmd/api

run-cli:
	go run ./cmd/cli

test:
	go test ./...

build:
	mkdir -p bin
	go build -o bin/api ./cmd/api
	go build -o bin/cli ./cmd/cli
