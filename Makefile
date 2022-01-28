.PHONY: default build build-static clean
default: build

SHELL	=	bash
out		=	ecurl

download:
	go get -d -v

build: download
	export GOOS=linux
	export GO111MODULE=on
	go build -o app

# Adds some flags for building the app statically linked
build-static: download
	export GOOS=linux
	export GO111MODULE=on
	export CGO_ENABLED=0
	go build \
		-ldflags="-extldflags=-static" \
		-tags osusergo,netgo \
		-o $(out)

clean:
	rm -rf ./$(out) ./vendor

test:
	go clean --testcache && go test ./... -v
