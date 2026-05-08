# Harbor Makefile

BINARY_NAME := harbor-cli
MAIN_PATH := cmd/harbor/main.go
VERSION_PKG := github.com/goharbor/harbor-cli/cmd/harbor/internal/version

# Extract version details from Git
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "none")
BUILD_DATE := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

# Build Configuration
GO_BUILD_ENV := CGO_ENABLED=0 GOOS=$(shell go env GOOS) GOARCH=$(shell go env GOARCH)

LDFLAGS := -s -w \
	-X '$(VERSION_PKG).Version=$(VERSION)' \
	-X '$(VERSION_PKG).Commit=$(COMMIT)' \
	-X '$(VERSION_PKG).Date=$(BUILD_DATE)'

.PHONY: all help build install clean test lint fmt tidy

all: help

## General
help: ## Display this help menu.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

## Development
fmt: ## format Go source code
	@echo "==> Formatting code"
	@go fmt ./...

lint: ## run golangci-lint for the whole project
	@echo "==> Running linter"
	@golangci-lint run ./...

tidy: ## clean up go.mod and go.sum
	@echo "==> Tidying modules"
	@go mod tidy

test: ## run unit tests with race detection and coverage
	@echo "==> Running tests"
	@go test -v -race -cover ./...

## Build cmd
harbor-cli: tidy ## build the harbor cli locally
	@$(GO_BUILD_ENV) go build -trimpath -ldflags "$(LDFLAGS)" -o $(BINARY_NAME) $(MAIN_PATH)
	@$(GO_BUILD_ENV) go install -trimpath -ldflags "$(LDFLAGS)" ./cmd/harbor
	@echo "==> Success creating ./$(BINARY_NAME) binary"

clean: ## remove the compiled binary and test cache
	@echo "==> Cleaning up"
	@rm -f $(BINARY_NAME)
	@go clean -testcache