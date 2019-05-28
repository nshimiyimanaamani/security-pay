.DEFAULT_GOAL = all

commit	:= $(shell git rev-parse --short HEAD)
goVersion	:= $(shell go version | cut -c12-19)
apiVersion	:= 0.0.1
os			:= $(shell go version | cut -c21-25)
arch		:= $(shell go version | cut -c27-31)


module:= $(shell basename "${PWD}")
package  := github.com/rugwirobaker/$(module)
packages := $(shell go list ./... | grep -v /vendor/)
image    := codebaker/$(module):$(commit)

.PHONY: all
all:: dependencies
all:: test build image

.PHONY: dependencies
dependencies::
	go mod tidy

.PHONY: build
build::
	CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -o \
	bin/app -ldflags="-w -s -X main.Service=${module} \
	-X main.GitCommit=${commit} \
	-X main.GoVersion=${goVersion} \
	-X main.GOOS=${os} \
	-X main.GOARCH=${arch}" main.go

.PHONY: test
test::
	go test -v $(packages)

.PHONY: bench
bench::
	go test -bench=. -v $(packages)

.PHONY: lint
lint::
	go vet -v $(packages)

.PHONY: check
check:: lint test

.PHONY: image
image::
	#building docker image ...
	docker build -t $(image) .

.PHONY: clean
clean::
	rm -f bin/$(module)