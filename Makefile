CWD=$(shell pwd)
GOPATH := $(CWD)

prep:
	if test -d pkg; then rm -rf pkg; fi

self:   prep
	if test -d src/github.com/whosonfirst/go-pubssed; then rm -rf src/github.com/whosonfirst/go-pubssed; fi
	mkdir -p src/github.com/whosonfirst/go-pubssed
	cp -r listener src/github.com/whosonfirst/go-pubssed/
	cp -r broker src/github.com/whosonfirst/go-pubssed/
	cp -r vendor/src/* src/

rmdeps:
	if test -d src; then rm -rf src; fi 

build:	fmt bin

deps:   rmdeps
	@GOPATH=$(GOPATH) go get -u "gopkg.in/redis.v1"
	@GOPATH=$(GOPATH) go get -u "github.com/facebookgo/grace/gracehttp"
	@GOPATH=$(GOPATH) go get -u "github.com/whosonfirst/go-httpony"

vendor-deps: deps
	if test ! -d vendor; then mkdir vendor; fi
	if test -d vendor/src; then rm -rf vendor/src; fi
	cp -r src vendor/src
	find vendor -name '.git' -print -type d -exec rm -rf {} +

bin: 	self
	@GOPATH=$(GOPATH) go build -o bin/pubssed-broadcast cmd/pubssed-broadcast.go
	@GOPATH=$(GOPATH) go build -o bin/pubssed-client cmd/pubssed-client.go
	@GOPATH=$(GOPATH) go build -o bin/pubssed-server cmd/pubssed-server.go

fmt:
	go fmt broker/*.go
	go fmt cmd/*.go
	go fmt listener/*.go
