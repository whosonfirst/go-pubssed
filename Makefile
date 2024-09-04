GOMOD=$(shell test -f "go.work" && echo "readonly" || echo "vendor")
LDFLAGS=-s -w

cli:
	go build -mod $(GOMOD) -ldflags="$(LDFLAGS)" -o bin/pubssed-broadcast cmd/pubssed-broadcast/main.go
	go build -mod $(GOMOD) -ldflags="$(LDFLAGS)" -o bin/pubssed-client cmd/pubssed-client/main.go
	go build -mod $(GOMOD) -ldflags="$(LDFLAGS)" -o bin/pubssed-server cmd/pubssed-server/main.go

docker:
	docker build -t go-pubssed .

install:
	sudo systemd/install.sh
