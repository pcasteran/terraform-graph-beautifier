all: build

setup:
	go get -u github.com/shurcooL/vfsgen

tools:
	go get -u golang.org/x/lint/golint

lint:
	golint ./...

fmt:
	go fmt ./...

generate:
	go generate -tags=dev ./...

build: generate
	go build -ldflags "-X main.version=$(git describe --tags)"

install:
	go install .

clean:
	go clean
