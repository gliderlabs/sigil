NAME=sigil
ARCH=$(shell uname -m)
ORG=gliderlabs
VERSION=0.6.0
define DESCRIPTION
Standalone string interpolator and template processor
Sigil is a command line tool for template processing
and POSIX-compliant variable expansion. It was created
for configuration templating, but can be used for any
text processing.
endef
REPO_NAME ?= gliderlabs/sigil
ARCHITECTURE = amd64

export DESCRIPTION

build:
	glu build darwin,linux ./cmd
ifdef CIRCLECI
	sudo chown -R circleci:circleci build
endif
	$(MAKE) deb
	$(MAKE) rpm

test:
	basht tests/*.bash

install: build
	install build/$(shell uname -s)/$(NAME) /usr/local/bin

deps:
	go get github.com/gliderlabs/glu
	go get -u github.com/progrium/basht/...
	go get -d ./cmd

release:
	glu release v$(VERSION)
	package_cloud push "dokku/dokku/ubuntu/xenial" build/deb/$(subst /,_,$(REPO_NAME))_$(VERSION)_amd64.deb
	package_cloud push "dokku/dokku/ubuntu/bionic" build/deb/$(subst /,_,$(REPO_NAME))_$(VERSION)_amd64.deb
	package_cloud push "dokku/dokku/ubuntu/focal" build/deb/$(subst /,_,$(REPO_NAME))_$(VERSION)_amd64.deb
	package_cloud push "dokku/dokku/debian/bullseye" build/deb/$(subst /,_,$(REPO_NAME))_$(VERSION)_amd64.deb
	package_cloud push "dokku/dokku/debian/buster" build/deb/$(subst /,_,$(REPO_NAME))_$(VERSION)_amd64.deb
	package_cloud push "dokku/dokku/debian/stretch" build/deb/$(subst /,_,$(REPO_NAME))_$(VERSION)_amd64.deb
	package_cloud push "dokku/dokku/el/7" build/rpm/$(subst /,-,$(REPO_NAME))-$(VERSION)-1.x86_64.rpm

clean:
	rm -rf build release

.PHONY: build release

deb:
	mkdir -p build/deb
	fpm -t deb -s dir -n $(subst /,-,$(REPO_NAME)) \
		 --version $(VERSION) \
		 --architecture amd64 \
		 --package build/deb/$(subst /,_,$(REPO_NAME))_$(VERSION)_amd64.deb \
		 --url "https://github.com/$(REPO_NAME)" \
		 --maintainer "Jose Diaz-Gonzalez <dokku@josediazgonzalez.com>" \
		 --category utils \
		 --description "$$DESCRIPTION" \
		 --license 'MIT License' \
		 build/Linux/sigil=/usr/bin/$(NAME)

rpm:
	mkdir -p build/rpm
	fpm -t rpm -s dir -n $(subst /,-,$(REPO_NAME)) \
		 --version $(VERSION) \
		 --architecture x86_64 \
		 --package build/rpm/$(subst /,-,$(REPO_NAME))-$(VERSION)-1.x86_64.rpm \
		 --url "https://github.com/$(REPO_NAME)" \
		 --category utils \
		 --description "$$DESCRIPTION" \
		 --license 'MIT License' \
		 build/Linux/sigil=/usr/bin/$(NAME)
