NAME=sigil
ARCH=$(shell uname -m)
ORG=gliderlabs
VERSION=0.3.2

build:
	glu build darwin,linux ./cmd

test: build
	basht tests/*.bash

install: build
	install build/$(shell uname -s)/sigil /usr/local/bin

deps:
	go get github.com/gliderlabs/glu
	go get -u github.com/progrium/basht/...
	go get ./cmd || true

release:
	glu release v$(VERSION)

clean:
	rm -rf build release

.PHONY: build release
