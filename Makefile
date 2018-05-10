PLUGINS_DIR=~/.terraform.d/plugins/darwin_amd64
PKG_NAME=dropbox

default: clean build move

build:
	go install

move:
	cp $(GOPATH)/bin/terraform-provider-$(PKG_NAME) $(PLUGINS_DIR)

clean:
	rm -f $(GOPATH)/bin/terraform-provider-$(PKG_NAME) $(PLUGINS_DIR)/terraform-provider-$(PKG_NAME)
