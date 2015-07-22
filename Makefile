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

deb:
	GOOS=linux GOARCH=amd64 go build -o usr/bin/puppet-anvil
	fpm -f -n puppet-anvil -s dir -t deb \
		--workdir debian \
		--version `git describe --tags --long` \
		--deb-upstart debian/upstart/puppet-anvil \
		--after-install debian/postinst usr/bin/
	rm -r usr

.PHONY: build


