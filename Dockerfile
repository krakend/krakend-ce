FROM debian:buster-slim

RUN apt-get update && apt-get install -y curl gcc
RUN curl https://storage.googleapis.com/golang/go1.17.4.linux-amd64.tar.gz | tar xzf - -C /usr/local

ENV GOPATH /go
ENV PATH $PATH:$GOPATH/bin:/usr/local/go/bin
ENV REPO_PATH $GOPATH/src/github.com/scriptdash/krakend-ce

ADD . $REPO_PATH
WORKDIR $REPO_PATH
RUN go install ./cmd/krakend-ce

FROM debian:buster-slim

COPY --from=0 /go/bin/krakend-ce /usr/bin/krakend

RUN apt-get update && \
	apt-get install -y ca-certificates && \
	update-ca-certificates && \
	rm -rf /var/lib/apt/lists/*

RUN useradd -r -c "KrakenD user" -U krakend

USER krakend

VOLUME [ "/etc/krakend" ]

WORKDIR /etc/krakend

ENTRYPOINT [ "/usr/bin/krakend" ]
CMD [ "run", "-c", "/etc/krakend/krakend.json" ]

EXPOSE 8000 8090
