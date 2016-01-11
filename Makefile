default: build

lint:
	golint ./...

test: lint
	godep go test ./...

build:
	godep go build

changelog:
	git log --first-parent --pretty="format:* %b" v`./banshee -v`..

.PHONY: changelog
