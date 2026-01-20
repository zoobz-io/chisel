.PHONY: test test-unit test-integration test-bench lint lint-fix coverage clean help check ci install-tools install-hooks build

.DEFAULT_GOAL := help

help: ## Display available commands
	@echo "chisel Development Commands"
	@echo "==========================="
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-18s\033[0m %s\n", $$1, $$2}'

test: ## Run all tests with race detector
	@go test -v -race ./... ./golang/... ./markdown/... ./typescript/... ./python/... ./rust/... ./testing/...

test-unit: ## Run unit tests only (short mode)
	@go test -v -race -short ./... ./golang/... ./markdown/... ./typescript/... ./python/... ./rust/... ./testing/...

test-integration: ## Run integration tests
	@go test -v -race ./testing/integration/...

test-bench: ## Run benchmarks
	@go test -bench=. -benchmem -benchtime=1s ./testing/benchmarks/...

lint: ## Run linters
	@golangci-lint run --config=.golangci.yml --timeout=5m ./... ./golang/... ./markdown/... ./typescript/... ./python/... ./rust/... ./testing/...

lint-fix: ## Run linters with auto-fix
	@golangci-lint run --config=.golangci.yml --fix ./... ./golang/... ./markdown/... ./typescript/... ./python/... ./rust/... ./testing/...

coverage: ## Generate coverage report (HTML)
	@go test -coverprofile=coverage.out ./... ./golang/... ./markdown/... ./typescript/... ./python/... ./rust/... ./testing/...
	@go tool cover -html=coverage.out -o coverage.html
	@go tool cover -func=coverage.out | tail -1
	@echo "Coverage report: coverage.html"

build: ## Build all packages
	@go build ./... ./golang/... ./markdown/... ./typescript/... ./python/... ./rust/... ./testing/...

clean: ## Remove generated files
	@rm -f coverage.out coverage.html coverage.txt
	@rm -rf dist/
	@find . -name "*.test" -delete
	@find . -name "*.prof" -delete
	@find . -name "*.out" -delete

install-tools: ## Install development tools
	@go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.7.2

install-hooks: ## Install git pre-commit hook
	@mkdir -p .git/hooks
	@echo '#!/bin/sh' > .git/hooks/pre-commit
	@echo 'make check' >> .git/hooks/pre-commit
	@chmod +x .git/hooks/pre-commit
	@echo "Pre-commit hook installed"

check: test lint ## Run tests and lint (quick validation)
	@echo "All checks passed!"

ci: clean lint test coverage test-bench ## Full CI simulation
	@echo "CI simulation complete!"
