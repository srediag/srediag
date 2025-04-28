# SPDX-License-Identifier: Apache-2.0
###############################################################################
# Project: srediag
# Description:
#   Most build logic is now handled by Mage targets (see magefiles/magefile.go).
#   The Makefile delegates to Mage for all major tasks.
###############################################################################

export GO111MODULE := on
export GOFLAGS      := -mod=readonly

.PHONY: all fmt lint test build build-plugins test-plugins deps install-dev-tools install-mage

all:
	mage all

fmt:
	mage format

lint:
	mage lint

test:
	mage test

build:
	mage build

build-plugins:
	mage buildplugins

test-plugins:
	mage testplugins

deps:
	mage deps

install-dev-tools:
	@echo ">> Installing development tools to $$(go env GOPATH)/bin..."
	@GOBIN=$$(go env GOPATH)/bin go install github.com/magefile/mage@v1.15.0
	@GOBIN=$$(go env GOPATH)/bin go install golang.org/x/tools/cmd/goimports@latest
	@GOBIN=$$(go env GOPATH)/bin go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	@GOBIN=$$(go env GOPATH)/bin go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	@curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin v2.1.5
	@echo "Development tools installed in $$(go env GOPATH)/bin ✓"

install-mage:
	@command -v mage >/dev/null 2>&1 || { \
	  echo ">> Installing Mage CLI…"; \
	  cd tools && go generate; \
	}
	@echo "Mage available at: $$(command -v mage)"
