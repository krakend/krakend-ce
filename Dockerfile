### BUILDER

FROM golang:1.12-alpine3.9 as BUILDER

RUN apk add --no-cache make curl git build-base

WORKDIR /build

COPY go.mod go.mod
COPY go.sum go.sum

RUN go mod download

COPY cmd cmd
COPY *.go ./

RUN go build cmd/krakend-ce/main.go
RUN chmod +x main

### RUNNER

FROM alpine:3.9

LABEL maintainer="dortiz@devops.faith"

RUN apk add --no-cache ca-certificates

COPY --from=BUILDER /build/main /usr/bin/krakend

VOLUME [ "/etc/krakend" ]

CMD [ "/usr/bin/krakend", "run", "-c", "/etc/krakend/krakend.json" ]

ENV GODEBUG=netdns=cgo 

EXPOSE 8000 8090
