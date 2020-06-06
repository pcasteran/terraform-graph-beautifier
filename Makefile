OUT := terraform-graph-beautifier
PKG := github.com/pcasteran/terraform-graph-beautifier
VERSION := $(shell git describe --always --long --dirty)

.PHONY: all
all: build

.PHONY: setup
setup:
	go get github.com/markbates/pkger/cmd/pkger

.PHONY: tools
tools:
	go get -u golang.org/x/lint/golint

.PHONY: lint
lint:
	golint ./...

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: tidy
tidy:
	go mod tidy

.PHONY: generate
generate:
	pkger

.PHONY: build
build: generate
	go build -i -v -o ${OUT}-v${VERSION} -ldflags="-X main.version=${VERSION}" ${PKG}

.PHONY: install
install:
	go install .

.PHONY: clean
clean:
	go clean

.PHONY: doc_generate
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