PLUGINS_DIR=~/.terraform.d/plugins/darwin_amd64
PKG_NAME=dropbox
VERSION=2.0.0

default: clean build

build:
	go build -o terraform-provider-$(PKG_NAME)_v$(VERSION)

move:
	cp ./terraform-provider-$(PKG_NAME)_v$(VERSION) $(PLUGINS_DIR)

clean:
	rm -f ./terraform-provider-$(PKG_NAME)_v$(VERSION) $(PLUGINS_DIR)/terraform-provider-$(PKG_NAME)_v$(VERSION)

.PHONY: default build move clean
