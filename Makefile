FILENAME:=$(shell ./filename.sh )

all: help

run-bench: build ## runs a new benchmark
	./run.sh | tee results/"$(FILENAME)"

charts: results/*.out  ## creates html charts from beanchmark output files in results/*
	for file in $^ ; do \
		./chart.py -f $${file} -o $${file}.html ; \
	done

build: ## compiles the whole code base
	@go version
	go build -v ./...

test: build ## executes all unit tests
	go clean -testcache
	go test ./...

fmt: ## uses gofmt to format the source code base
	gofmt -w $(shell find -name "*.go")

static-anal: ## executes basic static code-analysis tools
	staticcheck -f=stylish ./...
	go vet ./...
	go vet -vettool=$(shell which shadow) ./...

lint: ## runs a golang source code linter
	golint -set_exit_status ./...

help: ## display this help screen
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
