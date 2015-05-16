NAME=sigil
ARCH=$(shell uname -m)
VERSION=0.1.0

build:
	mkdir -p build/Linux  && GOOS=linux  go build -ldflags "-X main.Version $(VERSION)" -o build/Linux/$(NAME) ./cmd
	mkdir -p build/Darwin && GOOS=darwin go build -ldflags "-X main.Version $(VERSION)" -o build/Darwin/$(NAME) ./cmd

install: build
	install build/$(shell uname -s)/sigil /usr/local/bin

deps:
	go get -u github.com/progrium/gh-release/...
	go get || true

release:
	rm -rf release && mkdir release
	tar -zcf release/$(NAME)_$(VERSION)_Linux_$(ARCH).tgz -C build/Linux $(NAME)
	tar -zcf release/$(NAME)_$(VERSION)_Darwin_$(ARCH).tgz -C build/Darwin $(NAME)
	gh-release checksums sha256
	gh-release create gliderlabs/$(NAME) $(VERSION) $(shell git rev-parse --abbrev-ref HEAD) v$(VERSION)

clean:
	rm -rf build release

.PHONY: build release
