default: build

lint:
	golint ./...

test: lint
	go test ./...

build:
	go build

linux:
	GOOS=linux GOARCH=amd64 go build
