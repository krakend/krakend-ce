.PHONY: all build test

# This Makefile is a simple example that demonstrates usual steps to build a binary that can be run in the same
# architecture that was compiled in. The "ldflags" in the build assure that any needed dependency is included in the
# binary and no external dependencies are needed to run the service.

BIN_NAME :=krakend
OS := $(shell uname | tr '[:upper:]' '[:lower:]')
VERSION := 2.0.6
GIT_COMMIT := $(shell git rev-parse --short=7 HEAD)
PKGNAME := krakend
LICENSE := Apache 2.0
VENDOR=
URL := http://krakend.io
RELEASE := 0
USER := krakend
ARCH := amd64
DESC := High performance API gateway. Aggregate, filter, manipulate and add middlewares
MAINTAINER := Daniel Ortiz <dortiz@krakend.io>
DOCKER_WDIR := /tmp/fpm
DOCKER_FPM := devopsfaith/fpm
GOLANG_VERSION := 1.17.13
GLIBC_VERSION := $(shell sh find_glibc.sh)
ALPINE_VERSION := 3.15

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
	--depends ca-certificates \
	--before-remove builder/scripts/prerm.deb \
  --after-remove builder/scripts/postrm.deb \
	--before-install builder/scripts/preinst.deb

RPM_OPTS =--rpm-user $(USER) \
	--before-install builder/scripts/preinst.rpm \
	--before-remove builder/scripts/prerm.rpm \
  --after-remove builder/scripts/postrm.rpm

DEBNAME=${PKGNAME}_${VERSION}-${RELEASE}_${ARCH}.deb
RPMNAME=${PKGNAME}-${VERSION}-${RELEASE}.x86_64.rpm

all: test

build:
	@echo "Building the binary..."
	@go get .
	@go build -ldflags="-X github.com/luraproject/lura/v2/core.KrakendVersion=${VERSION} \
	-X github.com/luraproject/lura/v2/core.GoVersion=${GOLANG_VERSION} \
	-X github.com/luraproject/lura/v2/core.GlibcVersion=${GLIBC_VERSION}" \
	-o ${BIN_NAME} ./cmd/krakend-ce
	@echo "You can now use ./${BIN_NAME}"

test: build
	go test -v ./tests

# Build KrakenD using docker (defaults to whatever the golang container uses)
build_on_docker:
	docker run --rm -it -v "${PWD}:/app" -w /app golang:${GOLANG_VERSION} make -e build

# Build the container using the Dockerfile (alpine)
docker:
	docker build --no-cache --pull --build-arg GOLANG_VERSION=${GOLANG_VERSION} -t solamarpreet/krakend:${VERSION} .

benchmark:
	@mkdir -p bench_res
	@touch bench_res/${GIT_COMMIT}.out
	@docker run --rm -d --name krakend -v "${PWD}/tests/fixtures:/etc/krakend" -p 8080:8080 devopsfaith/krakend:${VERSION} run -dc /etc/krakend/bench.json
	@sleep 2
	@docker run --rm -it --link krakend peterevans/vegeta sh -c \
		"echo 'GET http://krakend:8080/test' | vegeta attack -rate=0 -duration=30s -max-workers=300 | tee results.bin | vegeta report" > bench_res/${GIT_COMMIT}.out
	@docker stop krakend
	@cat bench_res/${GIT_COMMIT}.out

security_scan:
	@mkdir -p sec_scan
	@touch sec_scan/${GIT_COMMIT}.out
	@docker run --rm -d --name krakend -v "${PWD}/tests/fixtures:/etc/krakend" -p 8080:8080 devopsfaith/krakend:${VERSION} run -dc /etc/krakend/bench.json
	@docker run --rm -it --link krakend instrumentisto/nmap --script vuln krakend > sec_scan/${GIT_COMMIT}.out
	@docker stop krakend
	@cat sec_scan/${GIT_COMMIT}.out

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

.PHONE: tgz
tgz: builder/skel/tgz/usr/bin/krakend
tgz: builder/skel/tgz/etc/krakend/krakend.json
tgz: builder/skel/tgz/etc/init.d/krakend
	tar zcvf krakend_${VERSION}_${ARCH}.tar.gz -C builder/skel/tgz/ .

.PHONY: deb
deb: builder/skel/deb/usr/bin/krakend
deb: builder/skel/deb/etc/krakend/krakend.json
	docker run --rm -it -v "${PWD}:${DOCKER_WDIR}" -w ${DOCKER_WDIR} ${DOCKER_FPM}:deb -t deb ${DEB_OPTS} \
		--iteration ${RELEASE} \
		--deb-systemd builder/files/krakend.service \
		-C builder/skel/deb \
		${FPM_OPTS}

.PHONY: rpm
rpm: builder/skel/rpm/usr/lib/systemd/system/krakend.service
rpm: builder/skel/rpm/usr/bin/krakend
rpm: builder/skel/rpm/etc/krakend/krakend.json
	docker run --rm -it -v "${PWD}:${DOCKER_WDIR}" -w ${DOCKER_WDIR} ${DOCKER_FPM}:rpm -t rpm ${RPM_OPTS} \
		--iteration ${RELEASE} \
		-C builder/skel/rpm \
		${FPM_OPTS}


.PHONY: clean
clean:
	rm -rf builder/skel/*
	rm -f *.deb
	rm -f *.rpm
	rm -f *.tar.gz
	rm -f krakend
	rm -rf vendor/
