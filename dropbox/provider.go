package dropbox

import (
	db "github.com/dropbox/dropbox-sdk-go-unofficial/dropbox"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

// Provider for the Dropbox API in Terraform
// returns a terraform.ResourceProvider
func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"access_token": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("DROPBOX_TOKEN", ""),
				Description: "Dropbox API access token generated from console",
			},
		},

		ResourcesMap: map[string]*schema.Resource{},

		DataSourcesMap: map[string]*schema.Resource{
			"dropbox_user_current": dataSourceDropboxUserCurrent(),
			"dropbox_space_usage":  dataSourceDropboxSpaceUsage(),
		},

		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	cfg := &ProviderConfig{
		DropboxConfig: &db.Config{Token: d.Get("token").(string)},
	}

	return cfg, nil
}
