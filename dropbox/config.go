package dropbox

import (
	db "github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ProviderConfig is passed in as meta data to all data
// source and resource management functions, containing
// the configuration struct for the Dropbox API
type ProviderConfig struct {
	DropboxConfig *db.Config
}

func providerConfig(d *schema.ResourceData) (interface{}, error) {
	cfg := &ProviderConfig{
		DropboxConfig: &db.Config{Token: d.Get("access_token").(string)},
	}

	return cfg, nil
}
