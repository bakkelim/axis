
# Makefile for building, testing, linting, and more

PROJECT_NAME := axis

.PHONY: build test lint format coverage

build:
	go build -o bin/$(PROJECT_NAME) ./...

test:
	go test ./... -v

lint:
	golangci-lint run

format:
	goimports -l -w .
	go fmt ./...

vet:
	go vet ./...

coverage:
	go test ./... -coverprofile=coverage.out
	go tool cover -func=coverage.out

precommit: format vet lint test
