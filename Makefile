# =============================================================================
# Configuration
# =============================================================================
GO ?= GOPRIVATE=github.com/motain go
VERSION ?= $(shell git rev-parse HEAD)
PACKAGES = $(shell go list -f {{.Dir}} ./... | grep -v /vendor/)
TOOLS_PATH := $(shell pwd)/tools

export PATH := ${TOOLS_PATH}:${PATH}
export GOBIN := ${TOOLS_PATH}

# =============================================================================
# Help
# =============================================================================
.PHONY: help
help: ## Show this help menu
	@echo ""
	@echo "Available targets:"
	@echo ""
	@grep -E '^[a-zA-Z0-9_/-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		sort | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}'
	@echo ""

# =============================================================================
# Build
# =============================================================================
.PHONY: build
build: ## Build the Linux binary
	@echo "Building Linux binary..."
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build \
		-ldflags="-X 'main.Version=${VERSION}'" \
		-o bin/linux/ofc \
		./cmd/${*}

# =============================================================================
# Code Quality
# =============================================================================
.PHONY: lint
lint: ## Run static analysis with staticcheck
	@echo "Running linter..."
	$(GO) install "honnef.co/go/tools/cmd/staticcheck@latest" && \
	$(GO) list -tags functional,unit ./... | \
		grep -v vendor/ | \
		grep -v /vendor/ | \
		grep -v /tools/ | \
		grep -v /mocks | \
		grep -v /wire_gen.go | \
		xargs -L1 staticcheck -tags functional,unit -f stylish -fail all -tests

.PHONY: test
test: ## Run all unit tests with coverage
	@echo "Running unit tests..."
	GOPRIVATE=github.com/motain $(GO) test -v -race -count=1 -tags unit ./... \
		-cover -coverprofile=coverage.out | \
		grep -v vendor/ | \
		grep -v /vendor/ | \
		grep -v /tools/ | \
		grep -v /mocks | \
		grep -v /wire_gen.go

.PHONY: stest
stest: ## Run tests for specific component (use C=path/to/component)
	@C=$${C:-""}; \
	C=$${C%/}; \
	echo "Running tests in component: ./$$C/..."; \
	GOPRIVATE=github.com/motain $(GO) test -race -count=1 -tags unit ./$$C/... \
		-cover -coverprofile=coverage.out | \
		grep -v vendor/ | \
		grep -v /vendor/ | \
		grep -v /tools/ | \
		grep -v /mocks/ | \
		grep -v /wire_gen.go && \
	cat coverage.out | grep -v "mocks" | grep -v "_gen.go" > cover.out

.PHONY: test/coverage
test/coverage: test ## Generate and display coverage statistics
	@echo "Generating coverage statistics..."
	$(GO) tool cover -func=coverage.out | \
		awk '/^[^total]/ {print $NF}' | \
		awk -F'%' '{sum+=$1; if(min==""){min=$1}; if($1>max){max=$1}; if($1<min){min=$1}; count++} END {print "Average:", sum/count "%"; print "Max:", max "%"; print "Min:", min "%"}'

# =============================================================================
# Dependencies
# =============================================================================
.PHONY: vendor
vendor: ## Vendor the dependencies
	@echo "Vendoring dependencies..."
	$(GO) mod tidy && $(GO) mod vendor && $(GO) mod verify

.PHONY: do-update-deps
do-update-deps: ## Update all dependencies
	@echo "Updating dependencies..."
	$(GO) get -u ./...

.PHONY: update-deps
update-deps: do-update-deps vendor ## Update dependencies and vendor them

# =============================================================================
# Code Generation
# =============================================================================
.PHONY: wire-all
wire-all: ## Run wire for all packages with wire.go files
	@echo "Running wire for all packages..."
	find . -type f -name wire.go -exec dirname {} \; | xargs wire gen

.PHONY: generate
generate: ## Generate code using go generate
	@echo "Generating code..."
	$(GO) generate ./...

# =============================================================================
# Cleanup
# =============================================================================
.PHONY: clean
clean: ## Remove service binaries
	@echo "Cleaning binaries..."
	rm -rf bin/

.PHONY: clean-state
clean-state: ## Clean the state directory (DESTRUCTIVE - requires confirmation)
	@echo "⚠️  WARNING: This will permanently delete all files in the state directory!"
	@echo "This operation cannot be undone."
	@echo ""
	@read -p "Are you sure you want to continue? (yes/no): " confirm && \
	if [ "$confirm" = "yes" ]; then \
		echo "Cleaning state directory..."; \
		rm -rf state/*; \
		echo "✅ State directory cleaned successfully."; \
	elif [ "$confirm" = "no" ]; then \
		echo "❌ Operation cancelled."; \
	else \
		echo "❌ Invalid input. Operation cancelled. Please answer 'yes' or 'no'."; \
		exit 1; \
	fi

# =============================================================================
# Application Management
# =============================================================================
.PHONY: create-metrics
create-metrics: ## Create metrics from configuration
	@echo "Creating metrics..."
	$(GO) run ./cmd/root.go metric apply -l ./config/grading-system/

.PHONY: create-scorecards
create-scorecards: ## Create scorecards from configuration
	@echo "Creating scorecards..."
	$(GO) run ./cmd/root.go scorecard apply -l ./config/scorecard/

.PHONY: create-components
create-components: ## Sync all components
	@echo "Syncing components..."
	./apply-all-individually.sh

.PHONY: bind-components
bind-components: ## Bind components to the grading system
	@echo "Binding components to grading system..."
	$(GO) run ./cmd/root.go component bind

.PHONY: create-all
create-all: ## Create all components, metrics, and scorecards (FULL SETUP)
	@echo "========================================"
	@echo "FULL SETUP - This will take a while!"
	@echo "========================================"
	@echo "Prerequisites:"
	@echo "  - Remove state directory if it exists"
	@echo "  - Ensure resources don't exist in compass"
	@echo "========================================"
	@echo "Starting full setup process..."
	$(MAKE) create-metrics
	$(MAKE) create-scorecards
	$(MAKE) create-components
	$(MAKE) bind-components
	@echo "========================================"
	@echo "All components, metrics, and scorecards created successfully!"
	@echo "========================================"