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

.PHONY: release
release:
	@if [ -z "$(version)" ]; then \
		echo ""; \
		echo "Error: version is not set. Please specify the version number."; \
		exit 1; \
	fi
	@git tag -a $(version) -m "Release $(version)"
	@git push origin $(version)
