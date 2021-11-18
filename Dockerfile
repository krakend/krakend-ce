FROM debian:buster-slim

LABEL maintainer="dortiz@devops.faith"

RUN apt-get update && \
	apt-get install -y ca-certificates && \
	update-ca-certificates && \
	rm -rf /var/lib/apt/lists/*

ADD krakend /usr/bin/krakend

ADD krakend.json /etc/krakend/

RUN useradd -M -u 1000 -c "KrakenD user" -U krakend

USER 1000

WORKDIR /etc/krakend

ENTRYPOINT [ "/usr/bin/krakend" ]
CMD [ "run", "-c", "/etc/krakend/krakend.json" ]

EXPOSE 8000 8090
