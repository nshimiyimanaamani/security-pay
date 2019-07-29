.DEFAULT_GOAL = all

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

.PHONY: dependencies
install::
	@echo "  >  installing dependencies..."
	@go mod tidy

.PHONY: build
build::
	@echo "  >  building binary..."
	@CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -o \
	bin/app -ldflags="-w -s -X ${package}/version.Service=${projectID}-backend \
	-X ${package}/version.GitCommit=${commit} \
	-X ${package}/version.GoVersion=${goVersion} \
	-X ${package}/version.main.GOOS=${os} \
	-X ${package}/version.main.GOARCH=${arch}" ./cmd/.

.PHONY: test
test::
	@echo "  >  running unit tests..."
	@go test -v $(packages)

.PHONY: bench
bench::
	@echo "  >  running benchmark tests..."
	@go test -bench=. -v $(packages)

.PHONY: lint
lint::
	@echo "  >  linting code..."
	@golint $(packages)

.PHONY: check
check:: lint test

.PHONY: coverage
coverage::
	@echo "  >  making coverage report..."
	@go test -cover $(packages)

.PHONY: image
image::
	@echo "  >  building docker image..."
	@docker build -t $(image) .

.PHONY:
push::
	@echo "  >  pushing docker image..."
	@docker tag $(image) $(registry)/$(projectID)/$(image)
	@docker push $(registry)/$(projectID)/$(image)
.PHONY: clean
clean::
	@echo "  >  cleaning..."
	@rm -f bin/$(module)