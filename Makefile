# SPDX-License-Identifier: Apache-2.0
###############################################################################
# Project: srediag
# Description:
#   Delegates tasks to Mage targets and enforces Go Modules (>=1.24).
###############################################################################

export GO111MODULE := on
export GOFLAGS      := -mod=readonly

GO_VERSION        := $(shell go version | awk '{print $$3}')
REQUIRED_GO_MAJOR := 1.24

# Get user home directory and GOPATH
HOME_DIR := $(shell echo $$HOME)
GOPATH := $(shell go env GOPATH)
GOBIN := $(GOPATH)/bin

# Output directories
PLUGIN_BASE_DIR := bin/plugins
PLUGIN_OUT_DIR := $(PLUGIN_BASE_DIR)
PLUGIN_TMP_DIR := $(PLUGIN_BASE_DIR)/.tmp

# Get Go version and build info
GO := go
GOARCH := $(shell go env GOARCH)
GOOS := $(shell go env GOOS)
VERSION := 0.1.0
BUILD_TIME := $(shell date -u '+%Y-%m-%d_%I:%M:%S%p')

# Protobuf variables
PROTO_FILES := $(shell find internal -name "*.proto")
PROTOC := protoc
PROTOC_GEN_GO := $(GOBIN)/protoc-gen-go
PROTOC_GEN_GO_GRPC := $(GOBIN)/protoc-gen-go-grpc
PROTOC_GEN_GOGOFAST := $(GOBIN)/protoc-gen-gogofast

# Build flags
BUILD_FLAGS := -trimpath \
    -ldflags="-X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME)"

# Plugin build flags
PLUGIN_BUILD_FLAGS := -trimpath \
    -tags=static \
    -installsuffix netgo \
    -ldflags="-X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME)" \
    $(BUILD_FLAGS)

# Binary paths
SREDIAG_BINARY := bin/srediag
SREDIAG_PLUGIN_BINARYS := bin/plugins/

# Configuration paths
CONFIG_DIR := configs
BUILDER_CONFIG := $(CONFIG_DIR)/srediag-builder.yaml
OTEL_CONFIG := $(CONFIG_DIR)/otel-config.yaml
TEST_CONFIG := $(CONFIG_DIR)/test-config.yaml
TEST_PLUGINS := $(CONFIG_DIR)/test-plugins.yaml

# Environment variables
export SREDIAG_CONFIG := $(abspath $(OTEL_CONFIG))
export SREDIAG_PLUGIN_DIR := $(abspath $(PLUGIN_OUT_DIR))

.PHONY: all check-env install-mage deps fmt lint test build clean docker install-dev-tools build-plugins update-deps verify-deps generate-plugins test-plugins generate-proto update-yaml-versions update-summary

# Build variables
BINARY_NAME := srediag
BUILD_DIR := bin
PLUGIN_DIR := $(BUILD_DIR)/plugins
PLUGIN_TMP_DIR := $(PLUGIN_DIR)/.tmp
PLUGIN_OUT_DIR := $(PLUGIN_DIR)

# Go environment
GO := go
REQUIRED_GO_MAJOR := 1.24
GO_VERSION := $(shell $(GO) version | cut -d" " -f3)
GOOS := $(shell go env GOOS)
GOARCH := $(shell go env GOARCH)

# Build flags
BUILD_FLAGS := -trimpath
PLUGIN_BUILD_FLAGS := -buildmode=plugin -trimpath

# Configuration files
BUILDER_CONFIG := configs/srediag-builder.yaml
TEST_CONFIG := configs/test-config.yaml
TEST_PLUGINS := configs/test-plugins.yaml

# Export plugin directory for use by the binary
export SREDIAG_PLUGIN_DIR := $(abspath $(PLUGIN_OUT_DIR))

.PHONY: all check-env install-mage deps fmt lint test build clean docker install-dev-tools build-plugins update-deps verify-deps generate-plugins test-plugins

all: check-env install-mage deps fmt lint test build test-plugins

install-dev-tools:
	@echo ">> Installing development tools to $(GOBIN)..."
	@mkdir -p $(GOBIN)
	@GOBIN=$(GOBIN) go install github.com/magefile/mage@v1.15.0
	@GOBIN=$(GOBIN) go install golang.org/x/tools/cmd/goimports@latest
	@GOBIN=$(GOBIN) go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	@GOBIN=$(GOBIN) go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	@curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(GOBIN) v2.1.5
	@echo "Development tools installed in $(GOBIN) ✓"

check-env:
	@echo ">> Checking Go environment…"
	@if ! expr "$(GO_VERSION)" : "go$(REQUIRED_GO_MAJOR)\." >/dev/null; then \
	  echo "Error: Go version must be $(REQUIRED_GO_MAJOR).x, found $(GO_VERSION)"; \
	  exit 1; \
	fi
	@echo "Go version $(GO_VERSION) ✓"
	@go env GO111MODULE | grep on >/dev/null 2>&1 || { \
	  echo "Error: GO111MODULE must be 'on'"; exit 1; }
	@echo "GO111MODULE=on ✓"

install-mage:
	@command -v mage >/dev/null 2>&1 || { \
	  echo ">> Installing Mage CLI…"; \
	  cd tools && go generate; \
	}
	@echo "Mage available at: $(shell command -v mage)"

update-deps:
	@echo ">> Updating dependencies..."
	@go get -u ./...
	@go mod tidy
	@echo "Dependencies updated ✓"

verify-deps:
	@echo ">> Verifying dependencies..."
	@go mod verify
	@echo "Dependencies verified ✓"

deps: verify-deps
	@echo ">> Downloading modules…"
	@go mod download

fmt:
	@echo ">> Formatting code…"
	mage Format

lint:
	@echo ">> Running linter…" 
	mage Lint

test:
	@echo ">> Running tests…"
	mage Test

build: build-binary generate-plugins build-plugins
	@echo ">> Building complete..."

build-binary:
	@echo ">> Building collector binary..."
	@echo "Using Go version: $(GO_VERSION)"
	@echo "Architecture: $(GOARCH)"
	@echo "OS: $(GOOS)"
	@mkdir -p bin
	@CGO_ENABLED=1 \
	GOARCH=$(GOARCH) \
	GOOS=$(GOOS) \
	$(GO) build \
		$(BUILD_FLAGS) \
		-o $(BUILD_DIR)/$(BINARY_NAME) \
		./cmd/srediag
	@echo "✓ Built collector binary at $(BUILD_DIR)/$(BINARY_NAME)"
	@if [ ! -x "$(BUILD_DIR)/$(BINARY_NAME)" ]; then \
		echo "Error: Failed to build collector binary"; \
		exit 1; \
	fi

generate-plugins: update-yaml-versions
	@echo ">> Generating plugins from $(BUILDER_CONFIG)..."
	@echo "Using Go version: $(GO_VERSION)"
	@echo "Architecture: $(GOARCH)"
	@echo "OS: $(GOOS)"
	@rm -rf $(PLUGIN_TMP_DIR)
	@mkdir -p $(PLUGIN_OUT_DIR)
	@mkdir -p $(PLUGIN_TMP_DIR)
	@SREDIAG_CONFIG=$(TEST_CONFIG) \
	 SREDIAG_PLUGIN_DIR=$(PLUGIN_TMP_DIR) \
	 $(BUILD_DIR)/$(BINARY_NAME) build generate \
	 --config $(BUILDER_CONFIG) \
	 --output-dir $(PLUGIN_TMP_DIR)
	@echo "✓ Generated plugins in $(PLUGIN_TMP_DIR)"

update-yaml-versions:
	./bin/srediag update-yaml-versions --yaml configs/srediag-builder.yaml --gomod go.mod --plugin-gen plugin/generated --makefile Makefile

# update-summary is auto-generated by the update utility
# Do not edit manually. To refresh, run `make update-yaml-versions` or `make`.

build-plugins:
	@echo ">> Building generated plugins..."
	@for plugin_dir in $(PLUGIN_TMP_DIR)/*; do \
		if [ -d "$$plugin_dir" ]; then \
			plugin_name=$$(basename $$plugin_dir); \
			plugin_type=$$(echo $$plugin_name | cut -d'_' -f1); \
			component_name=$$(echo $$plugin_name | cut -d'_' -f2-); \
			echo "Building $$plugin_type/$$component_name..."; \
			cd "$$plugin_dir" && \
			CGO_ENABLED=1 \
			GOARCH=$(GOARCH) \
			GOOS=$(GOOS) \
			$(GO) mod tidy && \
			$(GO) build \
				$(PLUGIN_BUILD_FLAGS) \
				-o ../../$$plugin_type/$$component_name \
				.; \
			build_status=$$?; \
			chmod +x ../../$$plugin_type/$$component_name; \
			cd - > /dev/null || exit 1; \
			if [ $$build_status -eq 0 ]; then \
				echo "✓ Successfully built $$plugin_type/$$component_name"; \
			else \
				echo "✗ Failed to build $$plugin_type/$$component_name"; \
				exit 1; \
			fi \
		fi \
	done
	@echo "✓ All plugins built successfully in $(PLUGIN_OUT_DIR)"

test-plugins: build
	@echo ">> Testing plugins..."
	@echo "Testing with configuration: $(TEST_CONFIG)"
	@for plugin_dir in $(PLUGIN_OUT_DIR)/*; do \
		if [ -d "$$plugin_dir" ]; then \
			plugin_type=$$(basename $$plugin_dir); \
			for plugin in $$plugin_dir/*; do \
				if [ -f "$$plugin" ]; then \
					plugin_name=$$(basename $$plugin); \
					echo "Testing $$plugin_type/$$plugin_name..."; \
					SREDIAG_CONFIG=$(TEST_CONFIG) \
					SREDIAG_PLUGIN_DIR=$(PLUGIN_OUT_DIR) \
					$(BUILD_DIR)/$(BINARY_NAME) build plugin $$plugin_type $$plugin_name \
						--config $(TEST_PLUGINS) || exit 1; \
					echo "✓ Plugin $$plugin_type/$$plugin_name tested successfully"; \
				fi \
			done \
		fi \
	done
	@echo "✓ All plugins tested successfully"

clean:
	@echo ">> Cleaning…"
	rm -rf $(BUILD_DIR)
	rm -rf $(PLUGIN_DIR)

# Docker targets
docker-build:
	@echo ">> Building Docker image…"
	docker build -t srediag/srediag:latest .

docker-run:
	docker run -p 4317:4317 -p 4318:4318 -p 55679:55679 srediag/srediag:latest

# Development helper targets
dev: deps build run

run: build
	./$(BUILD_DIR)/$(BINARY_NAME)

generate-proto:
	@echo ">> Generating protobuf code..."
	@for proto in $(PROTO_FILES); do \
		echo "Generating code for $$proto"; \
		$(PROTOC) -I. \
			--go_out=. --go_opt=paths=source_relative \
			--go-grpc_out=. --go-grpc_opt=paths=source_relative \
			$$proto; \
	done
	@echo "Protobuf code generated ✓"
