default: build

lint:
	(go list ./... | grep -v '/vendor/' )| while read -r line; do fgt golint "$$line" || exit 1; done

test: lint
	go test $$(go list ./... | grep -v '/vendor/')

build:
	go build

changelog:
	git log --first-parent --pretty="format:* %b" v`./banshee -v`..

static:
	make -C static deps
	make -C static build

.PHONY: changelog static
