NAME=sigil
ARCH=$(shell uname -m)
ORG=gliderlabs
VERSION=0.6.0

build:
	go build ./cmd/sigil

install: build
	install sigil /usr/local/bin/sigil

test:
	go test -v -race ./...

clean:
	rm -rf build release

.PHONY: build release
