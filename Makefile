MAKEFILE_LIST ?= Makefile

.PHONY: help
## Show this help message
help:
	@echo "Available targets:"
	@awk ' \
		BEGIN { help_seen = 0 } \
		/^##/ { sub(/^## ?/, "", $$0); help = $$0; help_seen = 1; next } \
		/^[a-zA-Z0-9._-]+:/ && help_seen { \
			split($$1, parts, ":"); \
			printf "  \033[36m%-12s\033[0m %s\n", parts[1], help; \
			help_seen = 0; \
		} \
	' $(MAKEFILE_LIST)

# Makefile
MAIN_NAME := $(shell basename $(shell git remote get-url origin) .git)
MAIN_PREFIX := $(MAIN_NAME)/
MAIN_VERSION := $(shell chmod +x ./scripts/version.sh 2>/dev/null; ./scripts/version.sh)

BINARY_NAME := $(MAIN_NAME)
OUTPUT_DIR := dist

.PHONY: fmt
## Format Go code and tidy go.mod
fmt:
	go mod tidy
	go fmt ./...

.PHONY: build
## Build the Go binary into the output directory
build:
	go build -o $(OUTPUT_DIR)/ \
		-ldflags "-X main.Name=$(MAIN_NAME) -X main.Version=$(MAIN_VERSION)"

.PHONY: test
## Run all Go tests
test:
	go test ./...
