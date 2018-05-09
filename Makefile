PKG_NAME=dropbox

default: clean build move

build:
	go build -o terraform-provider-$(PKG_NAME)

move:
	mv terraform-provider-$(PKG_NAME) ~/.terraform.d/plugins/darwin_amd64/

clean:
	rm ~/.terraform.d/plugins/darwin_amd64/terraform-provider-$(PKG_NAME)
