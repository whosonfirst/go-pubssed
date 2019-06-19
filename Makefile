tools:
	go build -mod vendor -o bin/pubssed-broadcast cmd/pubssed-broadcast/main.go
	go build -mod vendor -o bin/pubssed-client cmd/pubssed-client/main.go
	go build -mod vendor -o bin/pubssed-server cmd/pubssed-server/main.go

fmt:
	go fmt broker/*.go
	# go fmt cmd/*.go
	go fmt listener/*.go
