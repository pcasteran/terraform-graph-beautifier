PROJECT_NAME := "terraform-graph-beautifier"
PKG := "github.com/pcasteran/${PROJECT_NAME}"
VERSION := $(shell git describe --always --long --dirty)

.PHONY: all setup dep tidy lint fmt generate build dist install clean doc_generate

all: build

setup:
	go get -u golang.org/x/lint/golint && \
	go get -u github.com/markbates/pkger/cmd/pkger

dep:
	go mod download

tidy:
	go mod tidy

lint:
	golint -set_exit_status ./...

fmt:
	go fmt ./...

generate: dep
	pkger

build: dep generate
	go build -i -v -o ${PROJECT_NAME}-${VERSION} -ldflags="-X main.version=${VERSION}" ${PKG}

dist:
	goreleaser --snapshot --skip-publish --rm-dist

install:
	go install .

clean:
	go clean

doc_generate: install
	cd samples/config1/ && \
	terraform init && \
	\
	terraform graph | terraform-graph-beautifier \
		--exclude="module.root.provider" \
		--output-type=cyto-html \
		> ../../doc/config1.html && \
	\
	terraform graph | terraform-graph-beautifier \
		--exclude="module.root.provider" \
		--output-type=cyto-json \
		| jq . \
		> ../../doc/config1.json && \
	\
	terraform graph | terraform-graph-beautifier \
		--exclude="module.root.provider" \
		--output-type=graphviz \
		> ../../doc/config1.gv && \
	\
	cd -