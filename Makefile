PKG_NAME=dropbox

default: clean build

build:
	go install

clean:
	rm $(GOPATH)/bin/terraform-provider-$(PKG_NAME)
