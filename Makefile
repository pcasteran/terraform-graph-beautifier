OUT := terraform-graph-beautifier
PKG := github.com/pcasteran/terraform-graph-beautifier
VERSION := $(shell git describe --always --long --dirty)

all: build

setup:
	go get github.com/markbates/pkger/cmd/pkger

tools:
	go get -u golang.org/x/lint/golint

lint:
	golint ./...

fmt:
	go fmt ./...

tidy:
	go mod tidy

generate:
	pkger

build: generate
	go build -i -v -o ${OUT}-v${VERSION} -ldflags="-X main.version=${VERSION}" ${PKG}

install:
	go install .

clean:
	go clean
