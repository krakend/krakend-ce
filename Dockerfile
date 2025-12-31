ARG GOLANG_VERSION
ARG ALPINE_VERSION
FROM golang:${GOLANG_VERSION}-alpine${ALPINE_VERSION} as builder

RUN apk --no-cache --virtual .build-deps add make gcc musl-dev binutils-gold

COPY . /app
WORKDIR /app

# Build and validate plugin before building KrakenD
RUN make build-plugin
RUN make test-plugin

# Build KrakenD binary
RUN make build


FROM alpine:${ALPINE_VERSION}

LABEL maintainer="community@krakend.io"

RUN apk upgrade --no-cache --no-interactive && apk add --no-cache ca-certificates tzdata && \
    adduser -u 1000 -S -D -H krakend && \
    mkdir -p /etc/krakend/plugins && \
    echo '{ "version": 3 }' > /etc/krakend/krakend.json

COPY --from=builder /app/krakend /usr/bin/krakend
COPY --from=builder /app/plugins/static-content/hog-static-content.so /etc/krakend/plugins/

USER 1000

WORKDIR /etc/krakend

ENTRYPOINT [ "/usr/bin/krakend" ]
CMD [ "run", "-c", "/etc/krakend/krakend.json" ]

EXPOSE 8000 8090
