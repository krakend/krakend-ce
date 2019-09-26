.PHONY: all build test

# This Makefile is a simple example that demonstrates usual steps to build a binary that can be run in the same
# architecture that was compiled in. The "ldflags" in the build assure that any needed dependency is included in the
# binary and no external dependencies are needed to run the service.

BIN_NAME :=krakend
OS := $(shell uname | tr '[:upper:]' '[:lower:]')
VERSION := 1.0.0
PKGNAME := krakend
LICENSE := Apache 2.0
VENDOR=
URL := http://krakend.io
RELEASE := 0
USER := krakend
ARCH := amd64
DESC := High performance API gateway. Aggregate, filter, manipulate and add middlewares
MAINTAINER := Daniel Ortiz <dortiz@devops.faith>
DOCKER_WDIR := /tmp/fpm
DOCKER_FPM := devopsfaith/fpm
GOLANG_VERSION := 1.13

FPM_OPTS=-s dir -v $(VERSION) -n $(PKGNAME) \
  --license "$(LICENSE)" \
  --vendor "$(VENDOR)" \
  --maintainer "$(MAINTAINER)" \
  --architecture $(ARCH) \
  --url "$(URL)" \
  --description  "$(DESC)" \
	--config-files etc/ \
  --verbose

DEB_OPTS= -t deb --deb-user $(USER) \
	--before-remove builder/scripts/prerm.deb \
  --after-remove builder/scripts/postrm.deb \
	--before-install builder/scripts/preinst.deb

DEB_INIT=--deb-init builder/files/krakend.init

RPM_OPTS =--rpm-user $(USER) \
	--before-install builder/scripts/preinst.rpm \
	--before-remove builder/scripts/prerm.rpm \
  --after-remove builder/scripts/postrm.rpm

DEBNAME=${PKGNAME}_${VERSION}-${RELEASE}_${ARCH}.deb
RPMNAME=${PKGNAME}-${VERSION}-${RELEASE}.x86_64.rpm

all: test

build:
	@echo "Building the binary..."
	@GOPROXY=https://goproxy.io go get .
	@go build -ldflags="-X github.com/devopsfaith/krakend/core.KrakendVersion=${VERSION}" -o ${BIN_NAME} ./cmd/krakend-ce
	@echo "You can now use ./${BIN_NAME}"

test: build
	go test -v ./tests

docker_build:
	docker run --rm -it -v "${PWD}:/app" -w /app golang:${GOLANG_VERSION} make build

docker_build_alpine:
	docker build -t krakend_alpine_compiler builder/alpine
	docker run --rm -it -e "BIN_NAME=krakend-alpine" -v "${PWD}:/app" -w /app krakend_alpine_compiler make -e build

krakend_docker:
	@echo "You need to compile krakend using 'make docker_build_alpine' to build this container."
	docker build -t devopsfaith/krakend:${VERSION} .

tgz: builder/skel/tgz/usr/bin/krakend
tgz: builder/skel/tgz/etc/krakend/krakend.json
tgz: builder/skel/tgz/etc/init.d/krakend
	tar zcvf krakend_${VERSION}_${ARCH}.tar.gz -C builder/skel/tgz/ .

deb: ubuntu debian
rpm: el6 el7

ubuntu: ubuntu-trusty ubuntu-xenial
debian: debian-wheezy debian-jessie debian-stretch

builder/skel/el6/etc/init/krakend.conf: builder/files/krakend.conf.el6
	mkdir -p "$(dir $@)"
	cp builder/files/krakend.conf.el6 "$@"

builder/skel/%/etc/init/krakend.conf: builder/files/krakend.conf
	mkdir -p "$(dir $@)"
	cp builder/files/krakend.conf "$@"

builder/skel/%/etc/init.d/krakend: builder/files/krakend.init
	mkdir -p "$(dir $@)"
	cp builder/files/krakend.init "$@"

builder/skel/%/usr/bin/krakend: krakend
	mkdir -p "$(dir $@)"
	cp krakend "$@"

builder/skel/%/etc/krakend/krakend.json: krakend.json
	mkdir -p "$(dir $@)"
	cp krakend.json "$@"

builder/skel/%/lib/systemd/system/krakend.service: builder/files/krakend.service
	mkdir -p "$(dir $@)"
	cp builder/files/krakend.service "$@"

builder/skel/%/usr/lib/systemd/system/krakend.service: builder/files/krakend.service
	mkdir -p "$(dir $@)"
	cp builder/files/krakend.service "$@"

.PHONY: ubuntu-trusty
ubuntu-trusty: builder/skel/ubuntu-trusty/usr/bin/krakend
ubuntu-trusty: builder/skel/ubuntu-trusty/etc/krakend/krakend.json
ubuntu-trusty: builder/skel/ubuntu-trusty/etc/init.d/krakend
ubuntu-trusty: builder/skel/ubuntu-trusty/etc/init/krakend.conf
	docker run --rm -it -v "${PWD}:${DOCKER_WDIR}" -w ${DOCKER_WDIR} ${DOCKER_FPM}:deb -t deb ${DEB_OPTS} \
		--iteration ${RELEASE}.ubuntu-trusty \
		-C builder/skel/ubuntu-trusty \
		${DEB_INIT} \
		${FPM_OPTS}

.PHONY: ubuntu-xenial
ubuntu-xenial: builder/skel/ubuntu-xenial/usr/bin/krakend
ubuntu-xenial: builder/skel/ubuntu-xenial/etc/krakend/krakend.json
	docker run --rm -it -v "${PWD}:${DOCKER_WDIR}" -w ${DOCKER_WDIR} ${DOCKER_FPM}:deb -t deb ${DEB_OPTS} \
		--iteration ${RELEASE}.ubuntu-xenial \
		--deb-systemd builder/files/krakend.service \
		-C builder/skel/ubuntu-xenial \
		${FPM_OPTS}

.PHONY: debian-wheezy
debian-wheezy: builder/skel/debian-wheezy/usr/bin/krakend
debian-wheezy: builder/skel/debian-wheezy/etc/krakend/krakend.json
debian-wheezy: builder/skel/debian-wheezy/etc/init.d/krakend
debian-wheezy: builder/skel/debian-wheezy/etc/init/krakend.conf
	docker run --rm -it -v "${PWD}:${DOCKER_WDIR}" -w ${DOCKER_WDIR} ${DOCKER_FPM}:deb -t deb ${DEB_OPTS} \
		--iteration ${RELEASE}.debian-wheezy \
		-C builder/skel/debian-wheezy \
		--before-install builder/scripts/preinst-debian-wheezy.deb \
		${DEB_INIT} \
		${FPM_OPTS}

.PHONY: debian-jessie
debian-jessie: builder/skel/debian-jessie/usr/bin/krakend
debian-jessie: builder/skel/debian-jessie/etc/krakend/krakend.json
	docker run --rm -it -v "${PWD}:${DOCKER_WDIR}" -w ${DOCKER_WDIR} ${DOCKER_FPM}:deb -t deb ${DEB_OPTS} \
		--iteration ${RELEASE}.debian-jessie \
		--deb-systemd builder/files/krakend.service \
		-C builder/skel/debian-jessie \
		${FPM_OPTS}

.PHONY: debian-stretch
debian-stretch: builder/skel/debian-stretch/usr/bin/krakend
debian-stretch: builder/skel/debian-stretch/etc/krakend/krakend.json
	docker run --rm -it -v "${PWD}:${DOCKER_WDIR}" -w ${DOCKER_WDIR} ${DOCKER_FPM}:deb -t deb ${DEB_OPTS} \
		--iteration ${RELEASE}.debian-stretch \
		--deb-systemd builder/files/krakend.service \
		-C builder/skel/debian-stretch \
		${FPM_OPTS}

.PHONY: el7
el7: builder/skel/el7/usr/lib/systemd/system/krakend.service
el7: builder/skel/el7/usr/bin/krakend
el7: builder/skel/el7/etc/krakend/krakend.json
	docker run --rm -it -v "${PWD}:${DOCKER_WDIR}" -w ${DOCKER_WDIR} ${DOCKER_FPM}:rpm -t rpm ${RPM_OPTS} \
		--iteration ${RELEASE}.el7 \
		-C builder/skel/el7 \
		${FPM_OPTS}

.PHONY: el6
el6: builder/skel/el6/etc/init/krakend.conf
el6: builder/skel/el6/usr/bin/krakend
el6: builder/skel/el6/etc/krakend/krakend.json
	docker run --rm -it -v "${PWD}:${DOCKER_WDIR}" -w ${DOCKER_WDIR} ${DOCKER_FPM}:rpm -t rpm ${RPM_OPTS} \
		--iteration ${RELEASE}.el6 \
		-C builder/skel/el6 \
		${FPM_OPTS}

.PHONY: clean
clean:
	rm -rf builder/skel/*
	rm -f *.deb
	rm -f *.rpm
	rm -f *.tar.gz
	rm -f krakend
	rm -rf vendor/
