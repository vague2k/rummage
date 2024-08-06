# dont mind this makefile, I like pretty outputs :)
GREEN = \033[32m
BLUE = \033[34m
RESET = \033[0m

all: lint format test

.PHONY: lint
lint:
	@echo -e "$(BLUE)::$(RESET) Running golangci-lint";
	@output=$$(golangci-lint run --color=always); \
	if [ -z "$$output" ]; then \
	echo -e "$(GREEN)No linting issues found$(RESET)"; \
	else \
		echo -e "$$output"; \
		exit 1; \
	fi

.PHONY: format 
format:
	@echo -e "$(BLUE)::$(RESET) Running gofumpt formatter";
	@output=$$(gofumpt -l -w -d .); \
	if [ -z "$$output" ]; then \
	echo -e "$(GREEN)No formatting needed$(RESET)"; \
	else \
		echo -e "$$output"; \
	fi

.PHONY: test
test:
	@echo -e "$(BLUE)::$(RESET) Running tests";
# See https://github.com/gotestyourself/gotestsum
	@gotestsum --format-hide-empty-pkg; 

.PHONY: install
install:
	@echo -e "$(BLUE)::$(RESET) Installing local version of Rummage";
# See https://github.com/gotestyourself/gotestsum
	@go install
	@rummage -v

	
.PHONY: release
release:
	@if [ -z "$(version)" ]; then \
		echo ""; \
		echo "Error: version is not set. Please specify the version number."; \
		exit 1; \
	fi
	@git tag -a $(version) -m "Release $(version)"
	@git push origin $(version)
