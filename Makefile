SHELL = /bin/bash

.DEFAULT_GOAL := cli
.SUFFIXES =
.SUFFIXES = .go

GOCMD = go
BINARY_LINUX = gitfind
BINARY_WINDOWS = gitfind_win
ALL_FILES = ./...
BIN_DIR = ./bin

.PHONY:fmt
fmt:
	$(GOCMD) fmt -l -w $(ALL_FILES)

.PHONY:lint
lint: fmt
	$(GOCMD) lint $(ALL_FILES)

.PHONY:vet
vet: fmt
	$(GOCMD) vet $(ALL_FILES)

.PHONY:cli
cli: vet
	GOOS=linux GOARCH=amd64 \
			 $(GOCMD) build -v -o $(BIN_DIR)/$(BINARY_LINUX) ./cmd/cli/main.go
	GOOS=windows GOARCH=amd64 \
			 $(GOCMD) build -v -o $(BIN_DIR)/$(BINARY_WINDOWS) ./cmd/cli/main.go

.PHONY:test
test:
	$(GOCMD) test $(ALL_FILES)

.PHONY:clean
clean:
	$(GOCMD) clean
	rm -f ./$(BINARY_LINUX)

.PHONY:deps
deps:
	$(GOCMD) get

.PHONY:build-deps
build-deps:
	$(GOCMD) install golang.org/x/lint/golint@latest

