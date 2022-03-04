.PHONY: all build test

# This Makefile is a simple example that demonstrates usual steps to build a binary that can be run in the same
# architecture that was compiled in. The "ldflags" in the build assure that any needed dependency is included in the
# binary and no external dependencies are needed to run the service.

BIN_NAME :=krakend
OS := $(shell uname | tr '[:upper:]' '[:lower:]')
GIT_COMMIT := $(shell git rev-parse --short=7 HEAD)
VERSION := 2.0.0-alpha
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
GOLANG_VERSION := 1.17.3
ALPINE_VERSION := 3.14

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

update_krakend_deps:
	go get github.com/luraproject/lura/v2@v2.0.1
	go get github.com/devopsfaith/bloomfilter/v2@v2.0.0
	go get github.com/devopsfaith/krakend-amqp/v2@v2.0.0
	go get github.com/devopsfaith/krakend-botdetector/v2@v2.0.0
	go get github.com/devopsfaith/krakend-cel/v2@v2.0.0
	go get github.com/devopsfaith/krakend-circuitbreaker/v2@v2.0.0
	go get github.com/devopsfaith/krakend-cobra/v2@v2.0.0
	go get github.com/devopsfaith/krakend-cors/v2@v2.0.0
	go get github.com/devopsfaith/krakend-flexibleconfig/v2@v2.0.0
	go get github.com/devopsfaith/krakend-gelf/v2@v2.0.0
	go get github.com/devopsfaith/krakend-gologging/v2@v2.0.0
	go get github.com/devopsfaith/krakend-httpcache/v2@v2.0.0
	go get github.com/devopsfaith/krakend-httpsecure/v2@v2.0.0
	go get github.com/devopsfaith/krakend-influx/v2@v2.0.0
	go get github.com/devopsfaith/krakend-jose/v2@v2.0.0
	go get github.com/devopsfaith/krakend-jsonschema/v2@v2.0.0
	go get github.com/devopsfaith/krakend-lambda/v2@v2.0.0
	go get github.com/devopsfaith/krakend-logstash/v2@v2.0.0
	go get github.com/devopsfaith/krakend-lua/v2@v2.0.0
	go get github.com/devopsfaith/krakend-martian/v2@v2.0.0
	go get github.com/devopsfaith/krakend-metrics/v2@v2.0.0
	go get github.com/devopsfaith/krakend-oauth2-clientcredentials/v2@v2.0.0
	go get github.com/devopsfaith/krakend-opencensus/v2@v2.0.0
	go get github.com/devopsfaith/krakend-pubsub/v2@v2.0.0
	go get github.com/devopsfaith/krakend-ratelimit/v2@v2.0.0
	go get github.com/devopsfaith/krakend-rss/v2@v2.0.0
	go get github.com/devopsfaith/krakend-usage@v1.4.0
	go get github.com/devopsfaith/krakend-viper/v2@v2.0.0
	go get github.com/devopsfaith/krakend-xml/v2@v2.0.0
	make test

build:
	@echo "Building the binary..."
	@go get .
	@go build -ldflags="-X github.com/luraproject/lura/v2/core.KrakendVersion=${VERSION}" -o ${BIN_NAME} ./cmd/krakend-ce
	@echo "You can now use ./${BIN_NAME}"

test: build
	go test -v ./tests

build_on_docker:
	docker run --rm -it -v "${PWD}:/app" -w /app golang:${GOLANG_VERSION} make build

docker:
	docker build --pull --build-arg GOLANG_VERSION=${GOLANG_VERSION} --build-arg ALPINE_VERSION=${ALPINE_VERSION} -t devopsfaith/krakend:${VERSION} .

benchmark:
	@mkdir -p bench_res
	@touch bench_res/${GIT_COMMIT}.out
	@docker run --rm -d --name krakend -v "${PWD}/tests/fixtures:/etc/krakend" -p 8080:8080 devopsfaith/krakend:${VERSION} run -dc /etc/krakend/bench.json
	@sleep 2
	@docker run --rm -it --link krakend peterevans/vegeta sh -c \
		"echo 'GET http://krakend:8080/test' | vegeta attack -rate=0 -duration=30s -max-workers=300 | tee results.bin | vegeta report" > bench_res/${GIT_COMMIT}.out
	@docker stop krakend
	@cat bench_res/${GIT_COMMIT}.out

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
