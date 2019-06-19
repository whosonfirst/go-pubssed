# https://stackoverflow.com/questions/47028597/choosing-golang-docker-base-image
# https://hub.docker.com/_/golang/
FROM golang:1.12-alpine as builder

RUN mkdir /build

RUN apk update && apk upgrade \
    && apk add git make \
    && cd /build \       
    && git clone https://github.com/whosonfirst/go-pubssed.git \
    && cd go-pubssed \
    && go build -mod vendor -o /usr/bin/pubssed-server cmd/pubssed-server/main.go \
    && cd && rm -rf /build
