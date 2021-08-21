MAKEFLAGS += --always-make

.PHONY: help

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

%:
	@:

tests: ## Runs tests.
	@rm -rf coverage && mkdir -p coverage
	CGO_ENABLED=1 go test -mod=readonly -cover -covermode=atomic -coverprofile=coverage/profile.out .

benchmarks: ## Runs benchmarks.
	@clear
	go test -benchmem -bench=. -benchtime 10000x benchmark_test.go

coverage: tests ## Checks code coverage.
	go tool cover -func=coverage/profile.out
	go tool cover -html=coverage/profile.out -o coverage/report.html
