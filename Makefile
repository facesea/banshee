default: build

lint:
	golint ./...

test: lint
	godep go test ./...

build:
	godep go build
