ARG GOLANG_VERSION
FROM golang:${GOLANG_VERSION} as builder

COPY . /app
WORKDIR /app

RUN make build


FROM debian:stable-slim

RUN apt install -y ca-certificates && \
    adduser -u 1000 -S -D -H krakend && \
    mkdir /etc/krakend && \
    echo '{ "version": 3 }' > /etc/krakend/krakend.json

COPY --from=builder /app/krakend /usr/bin/krakend
COPY --from=builder /app/krakend.json /etc/krakend/krakend.json

USER 1000

WORKDIR /etc/krakend

ENTRYPOINT [ "/usr/bin/krakend" ]
CMD [ "run", "-c", "/etc/krakend/krakend.json" ]

EXPOSE 8080
