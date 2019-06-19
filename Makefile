tools:
	go build -mod vendor -o bin/pubssed-broadcast cmd/pubssed-broadcast/main.go
	go build -mod vendor -o bin/pubssed-client cmd/pubssed-client/main.go
	go build -mod vendor -o bin/pubssed-server cmd/pubssed-server/main.go

docker:
	docker build -t go-pubssed .
