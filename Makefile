all: lint format test

.PHONY: lint
lint:
	@echo "Running linter...";
	@golangci-lint run --color=always;

.PHONY: format
format:
	@echo "Running formatter...";
	@gofumpt -l -w -d .;

.PHONY: test
test:
	@echo "Running tests...";
	@set -euo pipefail
	@go test -json -v ./... 2>&1 | tee /tmp/gotest.log | gotestfmt -hide=empty-packages

