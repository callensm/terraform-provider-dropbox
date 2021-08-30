# Terraform Dropbox Provider

This Terraform provider was modeled off of the [unofficial Dropbox Golang SDK](https://github.com/dropbox/dropbox-sdk-go-unofficial). Follow the usage documentation from that repository to setup your Dropbox application and generate an access token prior to using this provider.

### Provider Usage

```hcl
provider "dropbox" {
  access_token = "${var.token}"
}
```

**access_token** is the generated auth token from Dropbox when you create an application through their developer console. It can be manually placed in the provider configuration as shown above, or if not provided, it will automatically attempt to retrieve the token value from the `DROPBOX_TOKEN` environment variable.
