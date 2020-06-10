.PHONY: help clean build test test-cover fmt coverage-out coverage-html

export GO111MODULE=on
export CGO_ENABLED=0

SRCS = $(shell git ls-files '*.go' | grep -v '^vendor/')

MAIN_DIRECTORY:=./cmd/tabia
BIN_OUTPUT:=$(if $(filter $(shell go env GOOS), Windows), bin/tabia.exe, bin/tabia)

TAG_NAME:=$(shell git tag -l --contains HEAD)
SHA:=$(shell git rev-parse HEAD)
VERSION:=$(if $(TAG_NAME),$(TAG_NAME),$(SHA))

help: ## Display this help message
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

clean: ## Clean build output
	@echo BIN_OUTPUT: ${BIN_OUTPUT}
	rm -rf bin/ cover.out

test: ## Run tests
	CGO_ENABLED=1 go test -v -race -count=1 ./...

test-cover: ## Run tests with coverage
	CGO_ENABLED=1 go test -v -race -count=1 -covermode=atomic -coverprofile=coverage.out ./...

build: clean ## Build binary
	@echo VERSION: $(VERSION)
	go build -v -trimpath -ldflags '-X "main.version=${VERSION}"' -o ${BIN_OUTPUT} ${MAIN_DIRECTORY}

outdated: ## Checks for outdated dependencies
	go list -u -m -json all | go-mod-outdated -update

fmt: ## formats all *.go files added to git
	gofmt -s -l -w $(SRCS)

coverage-out: test-cover ## Show coverage in cli
	@echo Coverage details
	@go tool cover -func=coverage.out

coverage-html: test-cover ## Show coverage in browser
	@go tool cover -html=coverage.out
