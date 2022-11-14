.PHONY: default clean build fmt lint vet cyclo ineffassign shellcheck errcheck goconst gosec abcgo style run test cover license before_commit help

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

build: ## Build binary containing service executable
	go build -ldflags="-X 'main.BuildTime=$(buildtime)' -X 'main.BuildVersion=$(version)' -X 'main.BuildBranch=$(branch)' -X 'main.BuildCommit=$(commit)'"

fmt: ## Run go fmt -w for all sources
	@echo "Running go formatting"
	./gofmt.sh

lint: ## Run golint
	@echo "Running go lint"
	./golint.sh

vet: ## Run go vet. Report likely mistakes in source code
	@echo "Running go vet"
	./govet.sh

cyclo: ## Run gocyclo
	@echo "Running gocyclo"
	./gocyclo.sh

ineffassign: ## Run ineffassign checker
	@echo "Running ineffassign checker"
	./ineffassign.sh

shellcheck: ## Run shellcheck
	shellcheck $(shell find . -name "*.sh")

errcheck: ## Run errcheck
	@echo "Running errcheck"
	./goerrcheck.sh

goconst: ## Run goconst checker
	@echo "Running goconst checker"
	./goconst.sh

gosec: ## Run gosec checker
	@echo "Running gosec checker"
	./gosec.sh

abcgo: ## Run ABC metrics checker
	@echo "Run ABC metrics checker"
	./abcgo.sh

style: fmt vet lint cyclo shellcheck errcheck goconst gosec ineffassign abcgo ## Run all the formatting related commands (fmt, vet, lint, cyclo) + check shell scripts

run: clean build ## Build the project and executes the binary
	./insights-results-aggregator-mock

test: clean build ## Run the unit tests
	@go test -coverprofile coverage.out $(shell go list ./... | grep -v tests)
	@go tool cover -func=coverage.out

cover: test ## Generate HTML pages with code coverage
	@go tool cover -html=coverage.out

coverage: ## Display code coverage on terminal
	@go tool cover -func=coverage.out

integration_tests: ## Run all integration tests
	@echo "Running all integration tests"
	@./test.sh

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
		awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'
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
