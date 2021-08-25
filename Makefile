PLUGINS_DIR=~/.terraform.d/plugins/darwin_amd64
PKG_NAME=dropbox
VERSION=3.0.0

.DEFAULT_GOAL := build

build: vendor
	go build -o terraform-provider-$(PKG_NAME)_v$(VERSION)

clean:
	rm -rf vendor/

test: vendor
	go test -v ./...

vendor: clean
	go mod tidy && go mod verify && go mod vendor

.PHONY: build clean test vendor
