.PHONY: default clean build fmt lint shellcheck abcgo style run test cover license before_commit help install_golangci-lint

SOURCES:=$(shell find . -name '*.go')
DOCFILES:=$(addprefix docs/packages/, $(addsuffix .html, $(basename ${SOURCES})))

version=0.1
branch=$(shell git rev-parse --abbrev-ref HEAD)
commit=$(shell git rev-parse HEAD)
buildtime=$(shell date)

default: build

clean: ## Run go clean
	@go clean
	rm -f rest-api-tests
	rm -rf coverage/
	rm -f covcounters*
	rm -f covmeta*
	rm -f coverage.txt

build: ## Build binary containing service executable
	go build -ldflags="-X 'main.BuildTime=$(buildtime)' -X 'main.BuildVersion=$(version)' -X 'main.BuildBranch=$(branch)' -X 'main.BuildCommit=$(commit)'"

build-cover:	${SOURCES}  ## Build binary with code coverage detection support
	go build -cover -ldflags="-X 'main.BuildTime=$(buildtime)' -X 'main.BuildVersion=$(version)' -X 'main.BuildBranch=$(branch)' -X 'main.BuildCommit=$(commit)'"

install_golangci-lint:
	go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest

fmt: install_golangci-lint ## Run go formatting
	@echo "Running go formatting"
	golangci-lint fmt

lint: install_golangci-lint ## Run go liting
	@echo "Running go linting"
	golangci-lint run --fix

shellcheck: ## Run shellcheck
	shellcheck $(shell find . -name "*.sh")

abcgo: ## Run ABC metrics checker
	@echo "Run ABC metrics checker"
	./abcgo.sh

style: fmt lint shellcheck abcgo ## Run all the formatting related commands (fmt, vet, lint, cyclo) + check shell scripts

run: clean build ## Build the project and executes the binary
	./insights-results-aggregator-mock

test: ## Run the unit tests
	@go test -coverprofile coverage.out $(shell go list ./... | grep -v tests)
	@go tool cover -func=coverage.out

cover: test ## Generate HTML pages with code coverage
	@go tool cover -html=coverage.out

coverage: ## Display code coverage on terminal
	@go tool cover -func=coverage.out

integration_tests: ## Run all integration tests
	@echo "Running all integration tests"
	@./test.sh

local_integration_tests: ## Run all integration tests locally
	@echo "Running all integration tests"
	@./rest-api-tests.sh

license:
	GO111MODULE=off go get -u github.com/google/addlicense && \
		addlicense -c "Red Hat, Inc" -l "apache" -v ./

before_commit: style test license
	./check_coverage.sh

help: ## Show this help screen
	@echo 'Usage: make <OPTIONS> ... <TARGETS>'
	@echo ''
	@echo 'Available targets are:'
	@echo ''
	@grep -E '^[ a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-25s\033[0m %s\n", $$1, $$2}'
	@echo ''

function_list: ${BINARY} ## List all functions in generated binary file
	go tool objdump ${BINARY} | grep ^TEXT | sed "s/^TEXT\s//g"

docs/packages/%.html: %.go
	mkdir -p $(dir $@)
	docgo -outdir $(dir $@) $^
	addlicense -c "Red Hat, Inc" -l "apache" -v $@

godoc: export GO111MODULE=off
godoc: install_docgo install_addlicense ${DOCFILES}

install_docgo: export GO111MODULE=off
install_docgo:
	[[ `command -v docgo` ]] || go get -u github.com/dhconnelly/docgo

install_addlicense: export GO111MODULE=off
install_addlicense:
	[[ `command -v addlicense` ]] || GO111MODULE=off go get -u github.com/google/addlicense
