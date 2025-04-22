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

.PHONY: all check-env install-mage deps fmt lint test build clean docker install-dev-tools

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

deps:
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

build:
	@echo ">> Building binary…"
	mage Build

clean:
	@echo ">> Cleaning…"
	rm -rf bin/

docker:
	@echo ">> Building Docker image…"
	docker build -t srediag/srediag:latest .
