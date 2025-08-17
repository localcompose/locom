# Makefile
MAIN_NAME := $(shell basename $(shell git remote get-url origin) .git)
MAIN_PREFIX := $(MAIN_NAME)/
MAIN_VERSION := $(shell \
  prefix=$$(echo $(MAIN_PREFIX)); \
  tag=$$(git describe --tags --abbrev=0  --match "$${prefix}*" 2>/dev/null || echo notags); \
  tagged_commit=$$(git rev-list -n 1 $$tag 2>/dev/null || echo ""); \
  current_commit=$$(git rev-parse HEAD); \
  short_commit=$$(git rev-parse --short HEAD); \
  version=$${tag#$$prefix}; \
  pre=""; \
  base=$$version; \
  if echo $$tag | grep -q '-'; then \
    pre=$${tag#*-}; \
    base=$${tag%-*}; \
  fi; \
  if [ "$$tagged_commit" != "$$current_commit" ]; then \
    if [ -n "$$pre" ]; then \
      base=$$base-$$pre.$$short_commit; \
    else \
      base=$$base-$$short_commit; \
    fi; \
  else \
    if [ -n "$$pre" ]; then \
      base=$$base-$$pre; \
    fi; \
  fi; \
  if [ -n "$$(git status --porcelain)" ]; then \
    ts=$$(date +%s); \
    commit_ts=$$(git log -1 --format=%ct); \
    extra=$$(($$ts - $$commit_ts)); \
    echo -n $$base+dev$$extra; \
  else \
    echo -n $$base; \
  fi; \
)

BINARY_NAME := $(MAIN_NAME)
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
	go fmt ./...

.PHONY: build
## Build the Go binary into the output directory
build:
	go build -o $(OUTPUT_DIR)/$(BINARY_NAME) \
		-ldflags "-X main.Name=$(MAIN_NAME) -X main.Version=$(MAIN_VERSION)"

.PHONY: test
## Run all Go tests
test:
	go test ./...
