default: build

test:
	go test -v ./...

build:
	go build

linux:
	GOOS=linux GOARCH=amd64 go build
