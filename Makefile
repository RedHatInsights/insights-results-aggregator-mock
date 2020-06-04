.PHONY: default clean build fmt lint vet cyclo ineffassign shellcheck errcheck goconst gosec abcgo style run test cover license before_commit help

SOURCES:=$(shell find . -name '*.go')

version=0.1
branch=$(shell git rev-parse --abbrev-ref HEAD)
commit=$(shell git rev-parse HEAD)
buildtime=$(shell date)

default: build

clean: ## Run go clean
	@go clean
	rm -f rest-api-tests

build: ## Run go build
	go build -ldflags="-X 'main.BuildTime=$(buildtime)' -X 'main.BuildVersion=$(version)' -X 'main.BuildBranch=$(branch)' -X 'main.BuildCommit=$(commit)'"

fmt: ## Run go fmt -w for all sources
	@echo "Running go formatting"
	gofmt -l .

lint: ## Run golint
	@echo "Running go lint"
	golint $(go list ./...)

vet: ## Run go vet. Report likely mistakes in source code
	@echo "Running go vet"
	go vet $(go list ./...)

cyclo: ## Run gocyclo
	@echo "Running gocyclo"
	gocyclo -over 9 -avg .

ineffassign: ## Run ineffassign checker
	@echo "Running ineffassign checker"
	ineffassign .

shellcheck: ## Run shellcheck
	shellcheck $(shell find . -name "*.sh")

errcheck: ## Run errcheck
	@echo "Running errcheck"
	errcheck ./...

goconst: ## Run goconst checker
	@echo "Running goconst checker"
	goconst -min-occurrences=3 ./...

gosec: ## Run gosec checker
	@echo "Running gosec checker"
	gosec ./...

abcgo: ## Run ABC metrics checker
	@echo "Run ABC metrics checker"
	abcgo -path .

style: fmt vet lint cyclo shellcheck errcheck goconst gosec ineffassign abcgo ## Run all the formatting related commands (fmt, vet, lint, cyclo) + check shell scripts

run: clean build ## Build the project and executes the binary
	./insights-content-service

test: clean build ## Run the unit tests
	@go test -coverprofile coverage.out $(shell go list ./... | grep -v tests)

cover: test
	@go tool cover -html=coverage.out

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
