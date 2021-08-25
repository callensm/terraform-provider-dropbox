package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"

	"github.com/callensm/terraform-provider-dropbox/dropbox"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: dropbox.Provider,
	})
}
