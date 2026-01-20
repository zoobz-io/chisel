.PHONY: help test lint check build clean install-tools

help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'

test: ## Run all tests
	go test -race -cover ./...

lint: ## Run linters
	golangci-lint run ./...

check: lint test ## Run linters and tests

build: ## Build all packages
	go build ./...

clean: ## Clean build artifacts
	rm -rf dist/ coverage.out coverage.html

install-tools: ## Install development tools
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

coverage: ## Generate coverage report
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
