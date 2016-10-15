CWD=$(shell pwd)
GOPATH := $(CWD)

prep:
	if test -d pkg; then rm -rf pkg; fi

self:   prep
	if test -d src/github.com/whosonfirst/go-httpony; then rm -rf src/github.com/whosonfirst/go-httpony; fi
	mkdir -p src/github.com/whosonfirst/go-httpony
	cp httpony.go src/github.com/whosonfirst/go-httpony/
	cp -r cors src/github.com/whosonfirst/go-httpony/
	cp -r tls src/github.com/whosonfirst/go-httpony/
	cp -r rewrite src/github.com/whosonfirst/go-httpony/
	cp -r crumb src/github.com/whosonfirst/go-httpony/
	cp -r crypto src/github.com/whosonfirst/go-httpony/
	cp -r sso src/github.com/whosonfirst/go-httpony/
	cp -r stats src/github.com/whosonfirst/go-httpony/
	cp -r vendor/src/* src/

rmdeps:
	if test -d src; then rm -rf src; fi 

deps:   
	@GOPATH=$(GOPATH) go get -u "github.com/vaughan0/go-ini"
	@GOPATH=$(GOPATH) go get -u "golang.org/x/net/html"
	@GOPATH=$(GOPATH) go get -u "golang.org/x/oauth2"

vendor-deps: rmdeps deps
	if test ! -d vendor; then mkdir vendor; fi
	if test -d vendor/src; then rm -rf vendor/src; fi
	cp -r src vendor/src
	find vendor -name '.git' -print -type d -exec rm -rf {} +

fmt:
	go fmt cmd/*.go
	go fmt *.go
	go fmt cors/*.go
	go fmt crumb/*.go
	go fmt rewrite/*.go
	go fmt sso/*.go
	go fmt stats/*.go
	go fmt tls/*.go

bin: 	self fmt
	@GOPATH=$(GOPATH) go build -o bin/echo-pony cmd/echo-pony.go
