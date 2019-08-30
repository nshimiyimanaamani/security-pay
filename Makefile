.DEFAULT_GOAL = help

commit	:= $(shell git rev-parse --short HEAD)
goVersion	:= $(shell go version | cut -c12-19)
os			:= $(shell go version | cut -c21-25)
arch		:= $(shell go version | cut -c27-31)


module		:= $(shell basename "${PWD}")
package  	:= github.com/rugwirobaker/$(module)
packages 	:= $(shell go list ./... | grep -v /vendor/)
registry 	:= gcr.io
projectID 	:= paypack
image    	:= $(module):$(commit)

all:: install
all:: test build image

build:   ## compile application binary into bin directory
	@echo "  >  building binary..."
	@CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -o \
	bin/app -ldflags="-w -s -X ${package}/version.Service=${projectID}-backend \
	-X ${package}/version.GitCommit=${commit} \
	-X ${package}/version.GoVersion=${goVersion} \
	-X ${package}/version.main.GOOS=${os} \
	-X ${package}/version.main.GOARCH=${arch}" ./cmd/.

test:  ## run unit tests
	@echo "  >  running unit tests..."
	@go test -v $(packages)

bench:  ## run benchmarks and generate report
	@echo "  >  running benchmark tests..."
	@go test -bench=. -v $(packages)



check: lint test ## assertain complaince

clean: ## clean out artefacts
	@echo "  >  cleaning..."
	@rm -f bin/$(module)

coverage::   ## generate test coverage report
	@echo ">  making test coverage report..."
	@go test -cover $(packages)

dev:  	      ## start development environment
	@echo "> starting dev environment..."
	@docker-compose up -d

dev-build:    ## start development environment
	@echo "> rebuilding dev environment..."
	@docker-compose up -d --build

dev-teardown: ## clean out development containers
	@echo "> cleaning dev environment..."
	@docker-compose down -v

image:
	@echo "  >  building docker image..."
	@docker build -t $(image) .

install:  ## install and verify package dependencies
	@echo "  >  installing dependencies..."
	@go mod tidy

lint:  # apply linting rules
	@echo "  >  linting code..."
	@golint $(packages)

push:
	@echo "  >  pushing docker image..."
	@docker tag $(image) $(registry)/$(projectID)/$(image)
	@docker push $(registry)/$(projectID)/$(image)


.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'