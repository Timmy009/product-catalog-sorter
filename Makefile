# Product Catalog Sorting System
# Enterprise-grade Makefile for Go project automation
# Author: Senior Software Engineer with 30+ years experience

.DEFAULT_GOAL := help
.PHONY: help build test lint clean coverage deps tidy run dev-setup ci benchmark profile docker

# Project configuration
PROJECT_NAME := product-catalog-sorting
BINARY_NAME := catalog-sorter
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT_HASH := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_TIME := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
GO_VERSION := $(shell go version | awk '{print $$3}')

# Directories
BUILD_DIR := bin
DIST_DIR := dist
COVERAGE_DIR := coverage
DOCS_DIR := docs

# Build flags
LDFLAGS := -ldflags "\
	-X main.Version=$(VERSION) \
	-X main.CommitHash=$(COMMIT_HASH) \
	-X main.BuildTime=$(BUILD_TIME) \
	-X main.GoVersion=$(GO_VERSION) \
	-w -s"

# Go build flags
BUILD_FLAGS := -trimpath -mod=readonly
RELEASE_FLAGS := $(BUILD_FLAGS) -a -installsuffix cgo

# Test flags
TEST_FLAGS := -race -timeout=30s
COVERAGE_FLAGS := -coverprofile=$(COVERAGE_DIR)/coverage.out -covermode=atomic

# Colors for output
RED := \033[0;31m
GREEN := \033[0;32m
YELLOW := \033[1;33m
BLUE := \033[0;34m
PURPLE := \033[0;35m
CYAN := \033[0;36m
NC := \033[0m # No Color

help: ## Display this help screen
	@echo "$(CYAN)$(PROJECT_NAME)$(NC) - Enterprise Product Catalog Sorting System"
	@echo ""
	@echo "$(YELLOW)Usage:$(NC)"
	@echo "  make [target]"
	@echo ""
	@echo "$(YELLOW)Available targets:$(NC)"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  $(GREEN)%-15s$(NC) %s\n", $$1, $$2}' $(MAKEFILE_LIST)
	@echo ""
	@echo "$(YELLOW)Build Information:$(NC)"
	@echo "  Version:     $(VERSION)"
	@echo "  Commit:      $(COMMIT_HASH)"
	@echo "  Go Version:  $(GO_VERSION)"
	@echo "  Build Time:  $(BUILD_TIME)"

## Development targets
dev-setup: deps ## Setup development environment
	@echo "$(BLUE)[INFO]$(NC) Setting up development environment..."
	@go install golang.org/x/tools/cmd/goimports@latest
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@go install github.com/air-verse/air@latest
	@mkdir -p $(BUILD_DIR) $(DIST_DIR) $(COVERAGE_DIR)
	@echo "$(GREEN)[SUCCESS]$(NC) Development environment ready"

deps: ## Download and verify dependencies
	@echo "$(BLUE)[INFO]$(NC) Downloading dependencies..."
	@go mod download
	@go mod verify
	@echo "$(GREEN)[SUCCESS]$(NC) Dependencies downloaded and verified"

tidy: ## Clean up dependencies
	@echo "$(BLUE)[INFO]$(NC) Tidying dependencies..."
	@go mod tidy
	@echo "$(GREEN)[SUCCESS]$(NC) Dependencies tidied"

## Build targets
build: ## Build the application for current platform
	@echo "$(BLUE)[INFO]$(NC) Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	@go build $(BUILD_FLAGS) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/main.go
	@echo "$(GREEN)[SUCCESS]$(NC) Built $(BUILD_DIR)/$(BINARY_NAME)"
	@ls -lh $(BUILD_DIR)/$(BINARY_NAME)

build-release: ## Build optimized release binary
	@echo "$(BLUE)[INFO]$(NC) Building release binary..."
	@mkdir -p $(BUILD_DIR)
	@CGO_ENABLED=0 go build $(RELEASE_FLAGS) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/main.go
	@echo "$(GREEN)[SUCCESS]$(NC) Release binary built"
	@ls -lh $(BUILD_DIR)/$(BINARY_NAME)

build-all: ## Build for all supported platforms
	@echo "$(BLUE)[INFO]$(NC) Building for all platforms..."
	@mkdir -p $(DIST_DIR)
	@for os in linux darwin windows; do \
		for arch in amd64 arm64; do \
			if [ "$$os" = "windows" ]; then ext=".exe"; else ext=""; fi; \
			echo "Building for $$os/$$arch..."; \
			GOOS=$$os GOARCH=$$arch CGO_ENABLED=0 go build \
				$(RELEASE_FLAGS) $(LDFLAGS) \
				-o $(DIST_DIR)/$(BINARY_NAME)-$$os-$$arch$$ext ./cmd/main.go; \
		done; \
	done
	@echo "$(GREEN)[SUCCESS]$(NC) Multi-platform build complete"
	@ls -lh $(DIST_DIR)/

## Test targets
test: ## Run all tests
	@echo "$(BLUE)[INFO]$(NC) Running all tests..."
	@go test $(TEST_FLAGS) ./...
	@echo "$(GREEN)[SUCCESS]$(NC) All tests passed"

test-unit: ## Run unit tests only
	@echo "$(BLUE)[INFO]$(NC) Running unit tests..."
	@go test $(TEST_FLAGS) ./internal/...
	@echo "$(GREEN)[SUCCESS]$(NC) Unit tests passed"

test-integration: ## Run integration tests only
	@echo "$(BLUE)[INFO]$(NC) Running integration tests..."
	@go test $(TEST_FLAGS) ./test/integration/...
	@echo "$(GREEN)[SUCCESS]$(NC) Integration tests passed"

test-verbose: ## Run tests with verbose output
	@echo "$(BLUE)[INFO]$(NC) Running tests with verbose output..."
	@go test $(TEST_FLAGS) -v ./...

coverage: ## Generate test coverage report
	@echo "$(BLUE)[INFO]$(NC) Generating coverage report..."
	@mkdir -p $(COVERAGE_DIR)
	@go test $(COVERAGE_FLAGS) ./...
	@go tool cover -html=$(COVERAGE_DIR)/coverage.out -o $(COVERAGE_DIR)/coverage.html
	@go tool cover -func=$(COVERAGE_DIR)/coverage.out | tail -1
	@echo "$(GREEN)[SUCCESS]$(NC) Coverage report: $(COVERAGE_DIR)/coverage.html"

benchmark: ## Run benchmarks
	@echo "$(BLUE)[INFO]$(NC) Running benchmarks..."
	@go test -bench=. -benchmem -run=^$$ ./...
	@echo "$(GREEN)[SUCCESS]$(NC) Benchmarks completed"

## Code quality targets
lint: ## Run linter
	@echo "$(BLUE)[INFO]$(NC) Running linter..."
	@golangci-lint run --config .golangci.yml
	@echo "$(GREEN)[SUCCESS]$(NC) Linting passed"

fmt: ## Format code
	@echo "$(BLUE)[INFO]$(NC) Formatting code..."
	@go fmt ./...
	@goimports -w .
	@echo "$(GREEN)[SUCCESS]$(NC) Code formatted"

vet: ## Run go vet
	@echo "$(BLUE)[INFO]$(NC) Running go vet..."
	@go vet ./...
	@echo "$(GREEN)[SUCCESS]$(NC) Vet analysis passed"

## Utility targets
clean: ## Clean build artifacts
	@echo "$(BLUE)[INFO]$(NC) Cleaning build artifacts..."
	@rm -rf $(BUILD_DIR) $(DIST_DIR) $(COVERAGE_DIR)
	@go clean -cache -testcache -modcache
	@echo "$(GREEN)[SUCCESS]$(NC) Clean completed"

run: build ## Build and run the application
	@echo "$(BLUE)[INFO]$(NC) Running application..."
	@./$(BUILD_DIR)/$(BINARY_NAME)

dev: ## Run in development mode with hot reload
	@echo "$(BLUE)[INFO]$(NC) Starting development server..."
	@air -c .air.toml

## Profiling targets
profile-cpu: ## Run CPU profiling
	@echo "$(BLUE)[INFO]$(NC) Running CPU profiling..."
	@go test -cpuprofile=cpu.prof -bench=. ./...
	@go tool pprof cpu.prof

profile-mem: ## Run memory profiling
	@echo "$(BLUE)[INFO]$(NC) Running memory profiling..."
	@go test -memprofile=mem.prof -bench=. ./...
	@go tool pprof mem.prof

## Docker targets
docker-build: ## Build Docker image
	@echo "$(BLUE)[INFO]$(NC) Building Docker image..."
	@docker build -t $(PROJECT_NAME):$(VERSION) -t $(PROJECT_NAME):latest .
	@echo "$(GREEN)[SUCCESS]$(NC) Docker image built"

docker-run: docker-build ## Run Docker container
	@echo "$(BLUE)[INFO]$(NC) Running Docker container..."
	@docker run --rm -it $(PROJECT_NAME):latest

## CI/CD targets
ci: deps lint vet test coverage ## Run full CI pipeline
	@echo "$(GREEN)[SUCCESS]$(NC) CI pipeline completed successfully"

pre-commit: fmt lint vet test ## Run pre-commit checks
	@echo "$(GREEN)[SUCCESS]$(NC) Pre-commit checks passed"

release: clean build-all test coverage ## Prepare release
	@echo "$(BLUE)[INFO]$(NC) Preparing release $(VERSION)..."
	@echo "$(GREEN)[SUCCESS]$(NC) Release $(VERSION) ready"

## Documentation targets
docs: ## Generate documentation
	@echo "$(BLUE)[INFO]$(NC) Generating documentation..."
	@go doc -all ./... > $(DOCS_DIR)/api.md
	@echo "$(GREEN)[SUCCESS]$(NC) Documentation generated"

## Security targets
security: ## Run security checks
	@echo "$(BLUE)[INFO]$(NC) Running security checks..."
	@go list -json -m all | nancy sleuth
	@gosec ./...
	@echo "$(GREEN)[SUCCESS]$(NC) Security checks passed"

## Performance targets
stress-test: ## Run stress tests
	@echo "$(BLUE)[INFO]$(NC) Running stress tests..."
	@go test -stress -timeout=5m ./...
	@echo "$(GREEN)[SUCCESS]$(NC) Stress tests completed"

load-test: build ## Run load tests
	@echo "$(BLUE)[INFO]$(NC) Running load tests..."
	@./scripts/load-test.sh
	@echo "$(GREEN)[SUCCESS]$(NC) Load tests completed"
