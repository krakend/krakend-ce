.PHONY: all prepare deps build

# This Makefile is a simple example that demonstrates usual steps to build a binary that can be run in the same
# architecture that was compiled in. The "ldflags" in the build assure that any needed dependency is included in the
# binary and no external dependencies are needed to run the service.

KRAKEND_VERSION=0.4.0
BIN_NAME=krakend
DEP_VERSION=0.3.2
OS := $(shell uname | tr '[:upper:]' '[:lower:]')

all: deps build

prepare:
	@echo "Installing dep..."
	@curl -L -s https://github.com/golang/dep/releases/download/v${DEP_VERSION}/dep-${OS}-amd64 -o ${GOPATH}/bin/dep
	@chmod a+x ${GOPATH}/bin/dep

deps:
	@echo "Setting up the vendors folder..."
	@dep ensure -v
	@echo ""
	@echo "Resolved dependencies:"
	@dep status
	@echo ""

build:
	@echo "Building the binary..."
	@go build -a -ldflags="-X github.com/devopsfaith/krakend/core.KrakendVersion=${KRAKEND_VERSION}" -o ${BIN_NAME}
	@echo "You can now use ./${BIN_NAME}"
