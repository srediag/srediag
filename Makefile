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

# Output directories
PLUGIN_OUT_DIR := bin/plugins
PLUGIN_TMP_DIR := $(PLUGIN_OUT_DIR)/.tmp

# Get Go version and build info
GO := go
GOARCH := $(shell go env GOARCH)
GOOS := $(shell go env GOOS)
VERSION := 0.1.0
BUILD_TIME := $(shell date -u '+%Y-%m-%d_%I:%M:%S%p')

# Build flags
BUILD_FLAGS := -trimpath \
    -ldflags="-X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME)"

# Plugin build flags - must match main binary flags
PLUGIN_BUILD_FLAGS := -buildmode=plugin \
    -trimpath \
    -tags=static \
    -installsuffix netgo \
    -ldflags="-linkmode external -extldflags '-static' -X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME)" \
    $(BUILD_FLAGS)

# Binary paths
SREDIAG_BINARY := bin/srediag

.PHONY: all check-env install-mage deps fmt lint test build clean docker install-dev-tools build-plugins update-deps verify-deps

all: check-env install-mage deps fmt lint test build

install-dev-tools:
	@echo ">> Installing development tools..."
	@go install github.com/magefile/mage@v1.15.0
	@go install golang.org/x/tools/cmd/goimports@latest
	@curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.54.2
	@echo "Development tools installed ✓"

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

build: build-binary build-plugins
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
		-o $(SREDIAG_BINARY) \
		./cmd/srediag
	@echo "✓ Built collector binary at $(SREDIAG_BINARY)"

build-plugins:
	@echo ">> Building plugins from otelcol-builder.yaml..."
	@echo "Using Go version: $(GO_VERSION)"
	@echo "Architecture: $(GOARCH)"
	@echo "OS: $(GOOS)"
	@rm -rf $(PLUGIN_TMP_DIR)
	@mkdir -p $(PLUGIN_OUT_DIR)
	@mkdir -p $(PLUGIN_TMP_DIR)
	@$(SREDIAG_BINARY) plugin generate --config otelcol-builder.yaml --output-dir $(PLUGIN_TMP_DIR)
	@echo ">> Building generated plugins..."
	@for plugin_dir in $(PLUGIN_TMP_DIR)/*; do \
		if [ -d "$$plugin_dir" ]; then \
			plugin_name=$$(basename $$plugin_dir); \
			plugin_type=$$(echo $$plugin_name | cut -d'_' -f1); \
			component_name=$$(echo $$plugin_name | cut -d'_' -f2-); \
			echo "Building $$component_name"; \
			mkdir -p $(PLUGIN_OUT_DIR)/$$plugin_type; \
			cd "$$plugin_dir" && \
			CGO_ENABLED=1 \
			GOARCH=$(GOARCH) \
			GOOS=$(GOOS) \
			$(GO) mod tidy && \
			$(GO) build \
				$(PLUGIN_BUILD_FLAGS) \
				-o $(CURDIR)/$(PLUGIN_OUT_DIR)/$$plugin_type/$${component_name}.so \
				.; \
			cd - > /dev/null || exit 1; \
		fi \
	done
	@echo "✓ Built plugins in $(PLUGIN_OUT_DIR)"

clean:
	@echo ">> Cleaning…"
	rm -rf bin/

docker:
	@echo ">> Building Docker image…"
	docker build -t srediag/srediag:latest .
