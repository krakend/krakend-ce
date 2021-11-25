ARG GOLANG_VERSION
ARG ALPINE_VERSION
FROM golang:${GOLANG_VERSION}-alpine${ALPINE_VERSION} as builder

RUN apk add make gcc musl-dev

COPY . /app
WORKDIR /app

RUN make build


FROM alpine:${ALPINE_VERSION}

LABEL maintainer="community@krakend.io"

RUN apk add --no-cache ca-certificates && \
    adduser -u 1000 -S -D -H krakend && \
    mkdir /etc/krakend && \
    echo '{ "version": 2 }' > /etc/krakend/krakend.json

COPY --from=builder /app/krakend /usr/bin/krakend

RUN useradd -M -u 1000 -c "KrakenD user" -U krakend

USER 1000

WORKDIR /etc/krakend

ENTRYPOINT [ "/usr/bin/krakend" ]
CMD [ "run", "-c", "/etc/krakend/krakend.json" ]

EXPOSE 8000 8090
