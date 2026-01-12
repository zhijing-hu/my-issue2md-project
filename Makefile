# Makefile for issue2md project

# Standardize build commands per constitution
.PHONY: all build test clean web

# Default target
all: build

# Build the main CLI application
build:
	go build -o bin/issue2md ./cmd/issue2md

# Run all tests
GOPATH?=$(shell go env GOPATH)
test:
	@echo "Running tests..."
	go test ./... -v

# Web service (future use per constitution reference)
web:
	@echo "Web service not yet implemented"

# Clean build artifacts
clean:
	rm -f bin/issue2md
	rm -f issue2md_*.log

# Install dependencies
.PHONY: deps
deps:
	go mod tidy

# Quick run for development
.PHONY: run
run: build
	./bin/issue2md "