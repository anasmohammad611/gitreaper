# Build configuration
BINARY_NAME=gitreaper
MAIN_PATH=./cmd/gitreaper
BUILD_DIR=build
VERSION?=dev

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod

.PHONY: all build clean test deps lint help

all: clean deps test build

build:
	$(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME) -v $(MAIN_PATH)

clean:
	$(GOCLEAN)
	rm -rf $(BUILD_DIR)

test:
	$(GOTEST) -v ./...

deps:
	$(GOMOD) tidy
	$(GOMOD) download

lint:
	golangci-lint run

help:
	@echo "Available targets:"
	@echo "  build    - Build the binary"
	@echo "  clean    - Clean build artifacts"
	@echo "  test     - Run tests"
	@echo "  deps     - Download dependencies"
	@echo "  lint     - Run linter"
	@echo "  help     - Show this help"