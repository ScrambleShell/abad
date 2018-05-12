.PHONY: all build test vendor

all: build test

build:
	go build -o ./cmd/abad/abad -v ./cmd/abad 

test:
	go test -race -v ./...

vendor:
	go get github.com/madlambda/vendor
	vendor
