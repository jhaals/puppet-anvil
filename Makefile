SHELL=/bin/bash

default: build

clean:
	rm -rf build tmp-test

deps:
	go get -u -t -v ./...

test:
	go test ./...

test-all: test acceptance-test

acceptance-test: build
	./acceptance-test.sh

build:
	mkdir -p build
	go build -o build/puppet-anvil

.PHONY: build


