default: build

test:
	go test ./...

build:
	go build

linux:
	GOOS=linux GOARCH=amd64 go build
