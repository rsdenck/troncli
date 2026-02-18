.PHONY: build release snapshot clean test

APP_NAME := troncli
VERSION := $(shell git describe --tags --always --dirty)
BUILD_DIR := bin

build:
	@echo "Building $(APP_NAME) $(VERSION)..."
	go build -ldflags "-s -w -X main.version=$(VERSION)" -trimpath -o $(BUILD_DIR)/$(APP_NAME) ./cmd/troncli

release:
	@echo "Creating release..."
	goreleaser release --clean

snapshot:
	@echo "Creating snapshot release..."
	goreleaser release --snapshot --clean

clean:
	@echo "Cleaning build directory..."
	rm -rf $(BUILD_DIR) dist

test:
	go test -v ./...
