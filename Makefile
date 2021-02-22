.PHONY: all build test

# This Makefile is a simple example that demonstrates usual steps to build a binary that can be run in the same
# architecture that was compiled in. The "ldflags" in the build assure that any needed dependency is included in the
# binary and no external dependencies are needed to run the service.

BIN_NAME :=krakend
OS := $(shell uname | tr '[:upper:]' '[:lower:]')
VERSION := 1.3.0
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
GOLANG_VERSION := 1.15.8

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

DEB_INIT=--deb-init builder/files/krakend.init

RPM_OPTS =--rpm-user $(USER) \
	--before-install builder/scripts/preinst.rpm \
	--before-remove builder/scripts/prerm.rpm \
  --after-remove builder/scripts/postrm.rpm

DEBNAME=${PKGNAME}_${VERSION}-${RELEASE}_${ARCH}.deb
RPMNAME=${PKGNAME}-${VERSION}-${RELEASE}.x86_64.rpm

all: test

update_krakend_deps:
	go get github.com/devopsfaith/krakend@master
	go get github.com/devopsfaith/krakend-amqp@master
	go get github.com/devopsfaith/krakend-botdetector@master
	go get github.com/devopsfaith/krakend-cel@master
	go get github.com/devopsfaith/krakend-circuitbreaker@master
	go get github.com/devopsfaith/krakend-cobra@master
	go get github.com/devopsfaith/krakend-consul@master
	go get github.com/devopsfaith/krakend-cors@master
	# go get github.com/devopsfaith/krakend-etcd@master
	go get github.com/devopsfaith/krakend-flexibleconfig@master
	go get github.com/devopsfaith/krakend-gelf@master
	go get github.com/devopsfaith/krakend-gologging@master
	go get github.com/devopsfaith/krakend-httpcache@master
	go get github.com/devopsfaith/krakend-httpsecure@master
	go get github.com/devopsfaith/krakend-jose@master
	go get github.com/devopsfaith/krakend-jsonschema@master
	go get github.com/devopsfaith/krakend-lambda@master
	go get github.com/devopsfaith/krakend-logstash@master
	go get github.com/devopsfaith/krakend-lua@master
	go get github.com/devopsfaith/krakend-martian@master
	go get github.com/devopsfaith/krakend-metrics@master
	go get github.com/devopsfaith/krakend-oauth2-clientcredentials@master
	go get github.com/devopsfaith/krakend-opencensus@master
	go get github.com/devopsfaith/krakend-pubsub@master
	go get github.com/devopsfaith/krakend-ratelimit@master
	go get github.com/devopsfaith/krakend-rss@master
	go get github.com/devopsfaith/krakend-usage@master
	go get github.com/devopsfaith/krakend-viper@master
	go get github.com/devopsfaith/krakend-xml@master
	make test

build:
	@echo "Building the binary..."
	@go get .
	@go build -ldflags="-X github.com/devopsfaith/krakend/core.KrakendVersion=${VERSION}" -o ${BIN_NAME} ./cmd/krakend-ce
	@echo "You can now use ./${BIN_NAME}"

test: build
	go test -v ./tests

build_on_docker:
	docker run --rm -it -v "${PWD}:/app" -w /app golang:${GOLANG_VERSION} make build

docker:
	docker build --pull -t devopsfaith/krakend:${VERSION} .

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
