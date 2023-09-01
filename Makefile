VERSION = "unset"
DATE=$(shell date -u +%Y-%m-%d-%H:%M:%S-%Z)


GIT_HASH=$(shell git rev-parse --short HEAD)
GOFILES=$(shell go list ./... | grep -v /vendor/)
IMAGE_DEV_TAG=dev
IMAGE_TAG:=tag
PROJECTNAME=$(shell basename "$(PWD)"
GOPROXY =$("https://proxy.golang.org")
BUILD_FLAGS = "-X github.com/nshimiyimanaamani/paypack-backend/pkg/build.version=$(VERSION) -X github.com/nshimiyimanaamani/paypack-backend/pkg/build.buildDate=$(DATE)"

all: help

.PHONY: build
build:  	## build development paypack binary
	@echo "> building all..."
	@CGO_ENABLED=0 go build -ldflags $(BUILD_FLAGS) -o bin/paypack ./cmd/paypack
	@CGO_ENABLED=0 go build -ldflags $(BUILD_FLAGS) -o bin/worker ./cmd/worker

build-main:
	@echo "> building main..."
	@CGO_ENABLED=0 go build -ldflags $(BUILD_FLAGS) -o bin/paypack ./cmd/paypack

build-worker:
	@echo "> building worker..."
	@CGO_ENABLED=0 go build -ldflags $(BUILD_FLAGS) -o bin/worker ./cmd/worker
	
clean:		## remove build artifacts
	@echo "> removing artifacts..."
	@rm -r bin/*

dev:  		## start development environment
	@echo "> starting dev environment..."
	@docker-compose up -d

run:  		## running development environment
	@echo "Running  paypack-umutekano..."
	@ ./bin/paypack

run-worker:  		## running development environment
	@echo "> running  paypack-umutekano-worker..."
	@ ./bin/worker

dev-build:  ## rebuild development environment
	@echo "> starting dev environment..."
	@docker-compose up -d --build
dev-teardown: ## clean up the development artifacts
	@echo "> cleaning dev environment..."
	@docker-compose down -v
image: 		## build docker image
	@echo "> building docker image..."
	@docker build -t docker.pkg.github.com/nshimiyimanaamani/paypack-backend/paypack:$(GIT_HASH) .
release:	## build the paypack server with version number
	@echo "> creating release binaries..."
	@CGO_ENABLED=0 GOOS=windows GOARCH=amd64 $(GOPROXY) go build -ldflags $(BUILD_FLAGS) -o bin/paypack_windows ./cmd/.
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOPROXY) go build -ldflags $(BUILD_FLAGS) -o bin/paypack_linux ./cmd/.
	@CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 $(GOPROXY) go build -ldflags $(BUILD_FLAGS) -o bin/paypack_darwin ./cmd/.

test:		## run unit tests
	@echo "> running unit tests..."
	@go test -race $(GOFILES)

tidy:		## verify dependencies
	@echo "> verifying dependincies..."
	@echo "> go mod tidy $(GOFILES)"

ui: ## compiling web assets into a dist folder
	@echo ">compiling web assets..."
	@npm install --prefix ./web && npm audit fix --prefix ./web --froce && npm run build --prefix ./web

.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
