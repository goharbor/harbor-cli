# Makefile for building the Harbor CLI project

# Variables
BINARY_NAME=harbor
BUILD_DIR=bin
SRC_DIR=cmd/harbor
MAIN_FILE=$(SRC_DIR)/main.go

# Retrieve the current Git commit hash
GIT_COMMIT := $(shell git rev-parse --short HEAD)

# Retrieve the current Git tag or fallback to describe
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "v0.0.1")

# Go version
GO_VERSION := $(shell go version | awk '{print $$3}')

# Build time
BUILD_TIME := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

# System information
SYSTEM := $(shell go env GOOS)/$(shell go env GOARCH)

# Output directory
OUTPUT_DIR := $(BUILD_DIR)

# ldflags for embedding version information
LDFLAGS := -X github.com/goharbor/harbor-cli/cmd/harbor/internal/version.Version=$(VERSION) \
           -X github.com/goharbor/harbor-cli/cmd/harbor/internal/version.GoVersion=$(GO_VERSION) \
           -X github.com/goharbor/harbor-cli/cmd/harbor/internal/version.GitCommit=$(GIT_COMMIT) \
           -X github.com/goharbor/harbor-cli/cmd/harbor/internal/version.BuildTime=$(BUILD_TIME) \
           -X github.com/goharbor/harbor-cli/cmd/harbor/internal/version.System=$(SYSTEM)

# Default target
.PHONY: all
all: build

# Build the binary with ldflags
.PHONY: build
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(OUTPUT_DIR)
	@echo "ldflags: $(LDFLAGS)"
	go build -ldflags "$(LDFLAGS)" -o $(OUTPUT_DIR)/$(BINARY_NAME) $(MAIN_FILE)
	@echo "Build complete: $(OUTPUT_DIR)/$(BINARY_NAME)"

# Clean the build artifacts
.PHONY: clean
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf $(BUILD_DIR)
	@echo "Clean complete."

# Display version information
.PHONY: version
version:
	@./$(OUTPUT_DIR)/$(BINARY_NAME) version

# Install dependencies
.PHONY: deps
deps:
	go mod download

# Run tests
.PHONY: test
test:
	@pushd ./test/e2e && go test -v && popd

# Help message
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  build    - Build the binary with version information"
	@echo "  clean    - Remove build artifacts"
	@echo "  version  - Display the version of the built binary"
	@echo "  deps     - Download and install dependencies"
	@echo "  test     - Run tests"
	@echo "  help     - Show this help message"
