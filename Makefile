GO ?= GOPRIVATE=github.com/motain go
VERSION ?=$(shell git rev-parse HEAD)
PACKAGES = $(shell go list -f {{.Dir}} ./... | grep -v /vendor/)
TOOLS_PATH := $(shell pwd)/tools

export PATH := ${TOOLS_PATH}:${PATH}
export GOBIN := ${TOOLS_PATH}

.PHONY: help
help: ## Show this help.
	@echo "Targets:"
	@grep -E '^[a-zA-Z\/_-]*:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\t%-20s: %s\n", $$1, $$2}'

.PHONY: build
build:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build \
		-ldflags="-X 'main.Version=${VERSION}'" \
		-o bin/linux/ofc \
		./cmd/${*}

.PHONY: lint
lint:
	$(GO) install "honnef.co/go/tools/cmd/staticcheck@latest" && $(GO) list -tags functional,unit ./...  | grep -v vendor/ | grep -v /vendor/ | grep -v /tools/ | grep -v /mocks | grep -v /wire_gen.go | xargs -L1 staticcheck -tags functional,unit -f stylish -fail all -tests

.PHONY: test
test:
	GOPRIVATE=github.com/motain $(GO) test -v -race -count=1 -tags unit ./... -cover -coverprofile=coverage.out | grep -v vendor/ | grep -v /vendor/ | grep -v /tools/ | grep -v /mocks | grep -v /wire_gen.go

.PHONY: stest
stest:
	@C=$${C:-""}; \
	C=$${C%/}; \
	echo "Running tests in C: ./$$C/..."; \
	GOPRIVATE=github.com/motain $(GO) test -race -count=1 -tags unit ./$$C/... -cover -coverprofile=coverage.out | grep -v vendor/ | grep -v /vendor/ | grep -v /tools/ | grep -v /mocks/ | grep -v /wire_gen.go && \
	cat coverage.out | grep -v "mocks" | grep -v "_gen.go" > cover.out

.PHONY: coverage/func
test/coverage: test
	$(GO) tool cover -func=coverage.out | awk '/^[^total]/ {print $NF}' | awk -F'%' '{sum+=$1; if(min==""){min=$1}; if($1>max){max=$1}; if($1<min){min=$1}; count++} END {print "Average:", sum/count "%"; print "Max:", max "%"; print "Min:", min "%"}'

.PHONY: vendor
vendor: ## Vendor the dependencies.
	$(GO) mod tidy && $(GO) mod vendor && $(GO) mod verify

.PHONY: do-update-deps
do-update-deps: ## Update dependencies.
	$(GO) get -u ./...

.PHONY: update-deps
update-deps: do-update-deps vendor ## Update dependencies and vendor.

.PHONY: clean
clean: ## Removes the service binary.
	rm -rf bin/

.PHONY: wire-all
wire-all:
	find . -type f -name wire.go -exec dirname {} \; | xargs wire gen

generate:
	$(GO) generate ./...
