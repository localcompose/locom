# Makefile

BINARY_NAME := locom
OUTPUT_DIR := dist

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

.PHONY: fmt
## Format Go code and tidy go.mod
fmt:
	go mod tidy
	go fmt

.PHONY: build
## Build the Go binary into the output directory
build:
	go build -o $(OUTPUT_DIR)/$(BINARY_NAME)

.PHONY: test
## Run all Go tests
test:
	go test ./...
