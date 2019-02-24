# Terraform Dropbox Provider

## Requirements

- Terraform 0.10.x or higher
- Go 1.8 or higher

## Building The Provider

Cloning the provider source code

```sh
$ mkdir -p $GOPATH/src/github.com/callensm; cd $GOPATH/src/github.com/callensm
$ git clone https://github.com/callensm/terraform-provider-dropbox.git
```

Build the source into executable binary

```sh
$ cd $GOPATH/src/github.com/callensm/terraform-provider-dropbox
$ make build
```

This will build the binary of the provider into the working directory for you to move to your Terraform plugins directory.

## Documentation

This Terraform provider was modeled off of the [unofficial Dropbox Golang SDK](https://github.com/dropbox/dropbox-sdk-go-unofficial). Follow the usage documentation from that repository to setup your Dropbox application and generate an access token prior to using this provider.

### Provider Usage

```hcl
provider "dropbox" {
  access_token = "${var.token}"
}
```

**access_token** is the generated auth token from Dropbox when you create an application through their developer console. It can be manually placed in the provider configuration as shown above, or if not provided, it will automatically attempt to retrieve the token value from the `DROPBOX_TOKEN` environment variable.
