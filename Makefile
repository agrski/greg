################################################################################
# Setup

SHELL = /bin/bash

.DEFAULT_GOAL := build
.SUFFIXES =
.SUFFIXES = .go

GOCMD = go
BINARY_LINUX = greg
BINARY_WINDOWS = greg.exe
ALL_FILES = ./...
BIN_DIR = ./bin

################################################################################
# Flags

ifneq ($(VERBOSE),)
	VERBOSE = -v
endif

################################################################################
# Targets

.PHONY:lint
lint:
	golangci-lint run --verbose

.PHONY:build
build: build-cli

.PHONY:build-cli
build-cli: clean lint
	GOOS=linux GOARCH=amd64 \
			 $(GOCMD) build -o $(BIN_DIR)/$(BINARY_LINUX) ./cmd/cli/
	GOOS=windows GOARCH=amd64 \
			 $(GOCMD) build -o $(BIN_DIR)/$(BINARY_WINDOWS) ./cmd/cli/

.PHONY:test
test: test-unit test-integration

.PHONY:test-unit
test-unit:
	$(GOCMD) test $(VERBOSE) $(ALL_FILES)

.PHONY:test-integration
test-integration:
	$(GOCMD) test $(VERBOSE) -tags integration $(ALL_FILES)

.PHONY:benchmark
benchmark:
	# Use -run to exclude non-benchmark tests
	$(GOCMD) test $(VERBOSE) -bench=. -benchmem -run=XXX ./pkg/...

.PHONY:clean
clean:
	$(GOCMD) clean
	rm -f ./bin/*

.PHONY:deps
deps:
	$(GOCMD) get

.PHONY:generate
generate:
	$(GOCMD) generate tools.go
