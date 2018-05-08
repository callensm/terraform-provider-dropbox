package main

import (
	"github.com/callensm/terraform-provider-dropbox/dropbox"
	"github.com/hashicorp/terraform/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: dropbox.Provider})
}
