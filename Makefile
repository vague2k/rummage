all: test release

.PHONY: all
all:
	$(MAKE) test && $(MAKE) release

.PHONY: test
test:
# this doesn't print any output for paths with no test files
	@go test ./... | grep -v '\[no test files\]'
# one of the tests, when ran locally, will "go get" a package, this cleans that up
	@go mod tidy

.PHONY: release
release:
	@if [ -z "$(version)" ]; then \
		echo ""; \
		echo "Error: version is not set. Please specify the version number."; \
		exit 1; \
	fi
	@git tag -a $(version) -m "Release $(version)"
	@git push origin $(version)

