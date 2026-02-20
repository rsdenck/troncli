# Makefile for TronCLI

VERSION := 0.2.19
BINARY_NAME := troncli

.PHONY: all clean build test install uninstall deb rpm aur snap docs

all: build

clean:
	rm -rf dist
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_NAME).exe

build:
	go build -ldflags="-s -w" -o $(BINARY_NAME) cmd/troncli/main.go

test:
	go test ./...

install:
	go install cmd/troncli/main.go

# Packaging Targets (require Linux environment)

deb:
	./scripts/build-deb.sh

rpm:
	./scripts/build-rpm.sh

aur:
	./scripts/build-aur.sh

snap:
	./scripts/build-snap.sh

# Documentation
docs:
	./scripts/generate-docs.sh
