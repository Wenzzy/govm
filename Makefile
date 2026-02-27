BINARY_NAME=govm
VERSION?=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME=$(shell date -u '+%Y-%m-%d_%H:%M:%S')
LDFLAGS=-ldflags "-s -w -X github.com/wenzzy/govm/internal/config.Version=$(VERSION) -X github.com/wenzzy/govm/internal/config.BuildTime=$(BUILD_TIME)"

GOOS?=$(shell go env GOOS)
GOARCH?=$(shell go env GOARCH)

.PHONY: all build install clean test lint release

all: build

build:
	@echo "Building $(BINARY_NAME)..."
	@go build $(LDFLAGS) -o bin/$(BINARY_NAME) ./cmd/govm

install: build
	@echo "Installing $(BINARY_NAME)..."
	@mkdir -p $(HOME)/.govm/bin
	@cp bin/$(BINARY_NAME) $(HOME)/.govm/bin/$(BINARY_NAME)
	@ln -sf $(HOME)/.govm/bin/$(BINARY_NAME) $(HOME)/.govm/bin/g
	@echo "Installed to $(HOME)/.govm/bin/$(BINARY_NAME)"
	@echo "Add to your PATH: export PATH=\"\$$HOME/.govm/bin:\$$HOME/.govm/current/bin:\$$PATH\""

uninstall:
	@echo "Uninstalling $(BINARY_NAME)..."
	@rm -f $(HOME)/.govm/bin/$(BINARY_NAME)
	@rm -f $(HOME)/.govm/bin/g

clean:
	@echo "Cleaning..."
	@rm -rf bin/
	@rm -rf dist/

test:
	@echo "Running tests..."
	@go test -v ./...

lint:
	@echo "Linting..."
	@golangci-lint run

# Build for all platforms
release: clean
	@echo "Building releases..."
	@mkdir -p dist
	@GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o dist/$(BINARY_NAME)-darwin-amd64 ./cmd/govm
	@GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o dist/$(BINARY_NAME)-darwin-arm64 ./cmd/govm
	@GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o dist/$(BINARY_NAME)-linux-amd64 ./cmd/govm
	@GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -o dist/$(BINARY_NAME)-linux-arm64 ./cmd/govm
	@echo "Release binaries built in dist/"

# Development helpers
dev: build
	@./bin/$(BINARY_NAME)

run:
	@go run ./cmd/govm $(ARGS)

deps:
	@go mod download
	@go mod tidy
	@go mod vendor

.DEFAULT_GOAL := build
