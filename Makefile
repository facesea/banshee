default: build

lint:
	golint ./...

test: lint
	godep go test ./...

build:
	godep go build

linux:
	GOOS=linux GOARCH=amd64 godep go build
