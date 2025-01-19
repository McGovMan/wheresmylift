FROM golang:1.23.4-bullseye AS builder

ARG VERSION
RUN test -n "$VERSION"

WORKDIR /go/src/app
COPY go.* ./
RUN go mod download

COPY . ./
RUN mkdir bin/ && go build -o bin/ -ldflags=-X=main.version=${VERSION}

FROM debian:bullseye-slim

RUN apt update && apt-get install -y ca-certificates curl && update-ca-certificates
COPY --from=builder /go/src/app/bin/api /go/bin/api

EXPOSE 80/tcp
ENTRYPOINT ["/go/bin/api"]
