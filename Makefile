.PHONY: build test deps install

build: test deps
	go build

deps:
	go build ./...

test: deps
	go test ./...

install: build
	go install