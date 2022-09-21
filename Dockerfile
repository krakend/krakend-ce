ARG GOLANG_VERSION
FROM golang:${GOLANG_VERSION}

COPY . /app
WORKDIR /app

RUN make build

RUN apt install -y ca-certificates && \
    adduser -u 1000 -S -D -H krakend && \
    mkdir /etc/krakend && \
    echo '{ "version": 3 }' > /etc/krakend/krakend.json

RUN cp /app/krakend /usr/bin/krakend

USER 1000

WORKDIR /etc/krakend

ENTRYPOINT [ "/usr/bin/krakend" ]
CMD [ "run", "-c", "/etc/krakend/krakend.json" ]

EXPOSE 8080
