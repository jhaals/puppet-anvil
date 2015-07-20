SHELL=/bin/bash

default: build

clean:
	rm -rf build tmp-test

deps:
	go get -u -t -v ./...

test:
	go test ./...

acceptance-test: build
	./acceptance-test.sh

build:
	mkdir -p build
	cd cli/server && go build -o ../../build/puppet-anvil

.PHONY: build


