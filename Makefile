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
	CGO_ENABLED=0 go build \
		-ldflags="-X 'main.Version=${VERSION}'" \
		-o bin/${*} \
		github.com/motain/of-catalog/${*}

.PHONY: lint
lint:
	$(GO) install "honnef.co/go/tools/cmd/staticcheck@latest" && $(GO) list -tags functional,unit ./...  | grep -v vendor/ | xargs -L1 staticcheck -tags functional,unit -f stylish -fail all -tests

.PHONY: test
test:
	$(GO) test -v -race -coverprofile=coverage.out -count=1 -tags unit ./... | grep -v vendor/

.PHONY: test/coverage
test/coverage: test
	$(GO) tool cover -html=coverage.out

.PHONY: test/functional
test/functional:
	$(GO) test -v -p=1 -count=1 -coverprofile=coverage.out -tags functional ./... | grep -v vendor/

.PHONY: trivy-docker
trivy-docker: build ## Builds a docker image for trivy vulnerability checks.
	docker build --no-cache -t trivy-test:test .

.PHONY: trivy
trivy: trivy-docker ## Runs trivy vulnerability checks
	trivy image trivy-test:test

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



