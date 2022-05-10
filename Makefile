NAME = gliderlabs-sigil
EMAIL = sigil@josediazgonzalez.com
MAINTAINER = gliderlabs
MAINTAINER_NAME = Jose Diaz-Gonzalez
REPOSITORY = sigil
HARDWARE = $(shell uname -m)
SYSTEM_NAME  = $(shell uname -s | tr '[:upper:]' '[:lower:]')
BASE_VERSION ?= 0.8.1
IMAGE_NAME ?= $(MAINTAINER)/$(REPOSITORY)
PACKAGECLOUD_REPOSITORY ?= dokku/dokku-betafish
BINARY_NAME = sigil

ifeq ($(CI_BRANCH),release)
	VERSION ?= $(BASE_VERSION)
	DOCKER_IMAGE_VERSION = $(VERSION)
else
	VERSION = $(shell echo "${BASE_VERSION}")build+$(shell git rev-parse --short HEAD)
	DOCKER_IMAGE_VERSION = $(shell echo "${BASE_VERSION}")build-$(shell git rev-parse --short HEAD)
endif

version:
	@echo "$(CI_BRANCH)"
	@echo "$(VERSION)"

define PACKAGE_DESCRIPTION
Standalone string interpolator and template processor
Sigil is a command line tool for template processing
and POSIX-compliant variable expansion. It was created
for configuration templating, but can be used for any
text processing.
endef

export PACKAGE_DESCRIPTION

LIST = build release release-packagecloud validate
targets = $(addsuffix -in-docker, $(LIST))

.env.docker:
	@rm -f .env.docker
	@touch .env.docker
	@echo "CI_BRANCH=$(CI_BRANCH)" >> .env.docker
	@echo "GITHUB_ACCESS_TOKEN=$(GITHUB_ACCESS_TOKEN)" >> .env.docker
	@echo "IMAGE_NAME=$(IMAGE_NAME)" >> .env.docker
	@echo "PACKAGECLOUD_REPOSITORY=$(PACKAGECLOUD_REPOSITORY)" >> .env.docker
	@echo "PACKAGECLOUD_TOKEN=$(PACKAGECLOUD_TOKEN)" >> .env.docker
	@echo "VERSION=$(VERSION)" >> .env.docker

build: prebuild
	@$(MAKE) build/darwin/$(NAME)
	@$(MAKE) build/linux/$(NAME)-amd64
	@$(MAKE) build/linux/$(NAME)-arm64
	@$(MAKE) build/linux/$(NAME)-armhf
	@$(MAKE) build/deb/$(NAME)_$(VERSION)_amd64.deb
	@$(MAKE) build/deb/$(NAME)_$(VERSION)_arm64.deb
	@$(MAKE) build/deb/$(NAME)_$(VERSION)_armhf.deb
	@$(MAKE) build/rpm/$(NAME)-$(VERSION)-1.x86_64.rpm

build-docker-image:
	docker build --rm -q -f Dockerfile.build -t $(IMAGE_NAME):build .

$(targets): %-in-docker: .env.docker
	docker run \
		--env-file .env.docker \
		--rm \
		--volume /var/lib/docker:/var/lib/docker \
		--volume /var/run/docker.sock:/var/run/docker.sock:ro \
		--volume /usr/bin/docker:/usr/local/bin/docker \
		--volume ${PWD}:/src/github.com/$(MAINTAINER)/$(REPOSITORY) \
		--workdir /src/github.com/$(MAINTAINER)/$(REPOSITORY) \
		$(IMAGE_NAME):build make -e $(@:-in-docker=)

build/darwin/$(NAME):
	mkdir -p build/darwin
	CGO_ENABLED=0 GOOS=darwin go build -a -asmflags=-trimpath=/src -gcflags=-trimpath=/src \
										-ldflags "-s -w -X main.Version=$(VERSION)" \
										-o build/darwin/$(NAME) cmd/sigil.go

build/linux/$(NAME)-amd64:
	mkdir -p build/linux
	CGO_ENABLED=0 GOOS=linux go build -a -asmflags=-trimpath=/src -gcflags=-trimpath=/src \
										-ldflags "-s -w -X main.Version=$(VERSION)" \
										-o build/linux/$(NAME)-amd64 cmd/sigil.go

build/linux/$(NAME)-arm64:
	mkdir -p build/linux
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -a -asmflags=-trimpath=/src -gcflags=-trimpath=/src \
										-ldflags "-s -w -X main.Version=$(VERSION)" \
										-o build/linux/$(NAME)-arm64 cmd/sigil.go
build/linux/$(NAME)-armhf:
	mkdir -p build/linux
	CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=5 go build -a -asmflags=-trimpath=/src -gcflags=-trimpath=/src \
										-ldflags "-s -w -X main.Version=$(VERSION)" \
										-o build/linux/$(NAME)-armhf cmd/sigil.go

build/deb/$(NAME)_$(VERSION)_amd64.deb: build/linux/$(NAME)-amd64
	export SOURCE_DATE_EPOCH=$(shell git log -1 --format=%ct) \
		&& mkdir -p build/deb \
		&& fpm \
		--architecture amd64 \
		--category utils \
		--description "$$PACKAGE_DESCRIPTION" \
		--input-type dir \
		--license 'MIT License' \
		--maintainer "$(MAINTAINER_NAME) <$(EMAIL)>" \
		--name $(NAME) \
		--output-type deb \
		--package build/deb/$(NAME)_$(VERSION)_amd64.deb \
		--url "https://github.com/$(MAINTAINER)/$(REPOSITORY)" \
		--vendor "" \
		--version $(VERSION) \
		--verbose \
		build/linux/$(NAME)-amd64=/usr/bin/$(BINARY_NAME) \
		LICENSE=/usr/share/doc/$(NAME)/copyright

build/deb/$(NAME)_$(VERSION)_arm64.deb: build/linux/$(NAME)-arm64
	export SOURCE_DATE_EPOCH=$(shell git log -1 --format=%ct) \
		&& mkdir -p build/deb \
		&& fpm \
		--architecture arm64 \
		--category utils \
		--description "$$PACKAGE_DESCRIPTION" \
		--input-type dir \
		--license 'MIT License' \
		--maintainer "$(MAINTAINER_NAME) <$(EMAIL)>" \
		--name $(NAME) \
		--output-type deb \
		--package build/deb/$(NAME)_$(VERSION)_arm64.deb \
		--url "https://github.com/$(MAINTAINER)/$(REPOSITORY)" \
		--vendor "" \
		--version $(VERSION) \
		--verbose \
		build/linux/$(NAME)-arm64=/usr/bin/$(BINARY_NAME) \
		LICENSE=/usr/share/doc/$(NAME)/copyright

build/deb/$(NAME)_$(VERSION)_armhf.deb: build/linux/$(NAME)-armhf
	export SOURCE_DATE_EPOCH=$(shell git log -1 --format=%ct) \
		&& mkdir -p build/deb \
		&& fpm \
		--architecture armhf \
		--category utils \
		--description "$$PACKAGE_DESCRIPTION" \
		--input-type dir \
		--license 'MIT License' \
		--maintainer "$(MAINTAINER_NAME) <$(EMAIL)>" \
		--name $(NAME) \
		--output-type deb \
		--package build/deb/$(NAME)_$(VERSION)_armhf.deb \
		--url "https://github.com/$(MAINTAINER)/$(REPOSITORY)" \
		--vendor "" \
		--version $(VERSION) \
		--verbose \
		build/linux/$(NAME)-armhf=/usr/bin/$(BINARY_NAME) \
		LICENSE=/usr/share/doc/$(NAME)/copyright

build/rpm/$(NAME)-$(VERSION)-1.x86_64.rpm: build/linux/$(NAME)-amd64
	export SOURCE_DATE_EPOCH=$(shell git log -1 --format=%ct) \
		&& mkdir -p build/rpm \
		&& fpm \
		--architecture x86_64 \
		--category utils \
		--description "$$PACKAGE_DESCRIPTION" \
		--input-type dir \
		--license 'MIT License' \
		--maintainer "$(MAINTAINER_NAME) <$(EMAIL)>" \
		--name $(NAME) \
		--output-type rpm \
		--package build/rpm/$(NAME)-$(VERSION)-1.x86_64.rpm \
		--rpm-os linux \
		--url "https://github.com/$(MAINTAINER)/$(REPOSITORY)" \
		--vendor "" \
		--version $(VERSION) \
		--verbose \
		build/linux/$(NAME)-amd64=/usr/bin/$(BINARY_NAME) \
		LICENSE=/usr/share/doc/$(NAME)/copyright

clean:
	rm -rf build release validation

ci-report:
	docker version
	rm -f ~/.gitconfig

docker-image:
	docker build --rm -q -f Dockerfile.hub -t $(IMAGE_NAME):$(DOCKER_IMAGE_VERSION) .

bin/gh-release:
	mkdir -p bin
	curl -o bin/gh-release.tgz -sL https://github.com/progrium/gh-release/releases/download/v2.3.3/gh-release_2.3.3_$(SYSTEM_NAME)_$(HARDWARE).tgz
	tar xf bin/gh-release.tgz -C bin
	chmod +x bin/gh-release

release: build bin/gh-release
	rm -rf release && mkdir release
	tar -zcf release/$(NAME)_$(VERSION)_linux_amd64.tgz -C build/linux $(NAME)-amd64
	tar -zcf release/$(NAME)_$(VERSION)_linux_arm64.tgz -C build/linux $(NAME)-arm64
	tar -zcf release/$(NAME)_$(VERSION)_linux_armhf.tgz -C build/linux $(NAME)-armhf
	tar -zcf release/$(NAME)_$(VERSION)_darwin_$(HARDWARE).tgz -C build/darwin $(NAME)
	cp build/deb/$(NAME)_$(VERSION)_amd64.deb release/$(NAME)_$(VERSION)_amd64.deb
	cp build/deb/$(NAME)_$(VERSION)_arm64.deb release/$(NAME)_$(VERSION)_arm64.deb
	cp build/deb/$(NAME)_$(VERSION)_armhf.deb release/$(NAME)_$(VERSION)_armhf.deb
	cp build/rpm/$(NAME)-$(VERSION)-1.x86_64.rpm release/$(NAME)-$(VERSION)-1.x86_64.rpm
	bin/gh-release create $(MAINTAINER)/$(REPOSITORY) $(VERSION) $(shell git rev-parse --abbrev-ref HEAD)

release-packagecloud:
	@$(MAKE) release-packagecloud-deb
	@$(MAKE) release-packagecloud-rpm

release-packagecloud-deb: build/deb/$(NAME)_$(VERSION)_amd64.deb build/deb/$(NAME)_$(VERSION)_arm64.deb build/deb/$(NAME)_$(VERSION)_armhf.deb
	package_cloud push $(PACKAGECLOUD_REPOSITORY)/ubuntu/bionic  build/deb/$(NAME)_$(VERSION)_amd64.deb
	package_cloud push $(PACKAGECLOUD_REPOSITORY)/ubuntu/focal   build/deb/$(NAME)_$(VERSION)_amd64.deb
	package_cloud push $(PACKAGECLOUD_REPOSITORY)/ubuntu/jammy   build/deb/$(NAME)_$(VERSION)_amd64.deb
	package_cloud push $(PACKAGECLOUD_REPOSITORY)/debian/stretch build/deb/$(NAME)_$(VERSION)_amd64.deb
	package_cloud push $(PACKAGECLOUD_REPOSITORY)/debian/buster  build/deb/$(NAME)_$(VERSION)_amd64.deb
	package_cloud push $(PACKAGECLOUD_REPOSITORY)/debian/bullseye build/deb/$(NAME)_$(VERSION)_amd64.deb
	package_cloud push $(PACKAGECLOUD_REPOSITORY)/ubuntu/focal    build/deb/$(NAME)_$(VERSION)_arm64.deb
	package_cloud push $(PACKAGECLOUD_REPOSITORY)/ubuntu/jammy    build/deb/$(NAME)_$(VERSION)_arm64.deb
	package_cloud push $(PACKAGECLOUD_REPOSITORY)/ubuntu/focal    build/deb/$(NAME)_$(VERSION)_armhf.deb
	package_cloud push $(PACKAGECLOUD_REPOSITORY)/ubuntu/jammy    build/deb/$(NAME)_$(VERSION)_armhf.deb
	package_cloud push $(PACKAGECLOUD_REPOSITORY)/raspbian/buster build/deb/$(NAME)_$(VERSION)_armhf.deb
	package_cloud push $(PACKAGECLOUD_REPOSITORY)/raspbian/bullseye build/deb/$(NAME)_$(VERSION)_armhf.deb

release-packagecloud-rpm: build/rpm/$(NAME)-$(VERSION)-1.x86_64.rpm
	package_cloud push $(PACKAGECLOUD_REPOSITORY)/el/7           build/rpm/$(NAME)-$(VERSION)-1.x86_64.rpm

validate:
	mkdir -p validation
	lintian build/deb/$(NAME)_$(VERSION)_amd64.deb || true
	lintian build/deb/$(NAME)_$(VERSION)_arm64.deb || true
	lintian build/deb/$(NAME)_$(VERSION)_armhf.deb || true
	dpkg-deb --info build/deb/$(NAME)_$(VERSION)_amd64.deb
	dpkg-deb --info build/deb/$(NAME)_$(VERSION)_arm64.deb
	dpkg-deb --info build/deb/$(NAME)_$(VERSION)_armhf.deb
	dpkg -c build/deb/$(NAME)_$(VERSION)_amd64.deb
	dpkg -c build/deb/$(NAME)_$(VERSION)_arm64.deb
	dpkg -c build/deb/$(NAME)_$(VERSION)_armhf.deb
	cd validation && ar -x ../build/deb/$(NAME)_$(VERSION)_amd64.deb
	cd validation && ar -x ../build/deb/$(NAME)_$(VERSION)_arm64.deb
	cd validation && ar -x ../build/deb/$(NAME)_$(VERSION)_armhf.deb
	cd validation && rpm2cpio ../build/rpm/$(NAME)-$(VERSION)-1.x86_64.rpm > $(NAME)-$(VERSION)-1.x86_64.cpio
	ls -lah build/deb build/rpm validation
	sha1sum build/deb/$(NAME)_$(VERSION)_amd64.deb
	sha1sum build/deb/$(NAME)_$(VERSION)_arm64.deb
	sha1sum build/deb/$(NAME)_$(VERSION)_armhf.deb
	sha1sum build/rpm/$(NAME)-$(VERSION)-1.x86_64.rpm
	go get -u github.com/progrium/basht/...
	basht tests/*.bash

prebuild:
	true
