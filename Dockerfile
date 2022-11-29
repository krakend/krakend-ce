ARG GOLANG_VERSION
ARG ALPINE_VERSION
FROM golang:${GOLANG_VERSION}-alpine${ALPINE_VERSION} as builder
ARG GITHUB_TOKEN

RUN apk --no-cache --virtual .build-deps add make gcc musl-dev binutils-gold git

COPY . /app
WORKDIR /app

ENV GITHUB_TOKEN=${GITHUB_TOKEN}
ENV GOPRIVATE=github.com/optivainc/* 
RUN git config --global url."https://${GITHUB_TOKEN}@github.com/optivainc".insteadOf "https://github.com/optivainc"
RUN make build


FROM alpine:${ALPINE_VERSION}

LABEL maintainer="community@krakend.io"

RUN apk add --no-cache ca-certificates && \
    adduser -u 1000 -S -D -H krakend && \
    mkdir /etc/krakend && \
    echo '{ "version": 3 }' > /etc/krakend/krakend.json

COPY --from=builder /app/krakend /usr/bin/krakend

USER 1000

WORKDIR /etc/krakend

ENTRYPOINT [ "/usr/bin/krakend" ]
CMD [ "run", "-c", "/etc/krakend/krakend.json" ]

EXPOSE 8000 8090
