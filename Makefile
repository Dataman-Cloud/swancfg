PACKAGES = $(shell go list ./...)
TEST_PACKAGES = $(shell go list ./... | grep -v vendor)

.PHONY: build fmt test test-cover-html test-cover-func collect-cover-data

# Prepend our vendor directory to the system GOPATH
# so that import path resolution will prioritize
# our third party snapshots.
export GO15VENDOREXPERIMENT=1
# GOPATH := ${PWD}/vendor:${GOPATH}
# export GOPATH

default: build

build: fmt build

build:
	go build -v -o swancfg main.go 

install:
	install -v swancfg /usr/local/bin

clean:
	rm -f swancfg
fmt:
	go fmt ./...

test:
	go test -cover=true ${TEST_PACKAGES}

collect-cover-data:
	@echo "mode: count" > coverage-all.out
	$(foreach pkg,$(TEST_PACKAGES),\
                go test -v -coverprofile=coverage.out -covermode=count $(pkg) || exit $?;\

