# Terraform Dropbox Provider

## Requirements

* Terraform 0.10.x or higher
* Go 1.8 or higher

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

## Documentation

This Terraform provider was modeled off of the [unofficial Dropbox Golang SDK](https://github.com/dropbox/dropbox-sdk-go-unofficial). Follow the usage documentation from that repository to setup your Dropbox application and generate an access token prior to using this provider.

Data Source [documentation](.github/DATA_SOURCES.md)

Resource [documentation](.github/RESOURCES.md)
