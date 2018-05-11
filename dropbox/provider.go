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
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("DROPBOX_TOKEN", ""),
				Description: "Dropbox API access token generated from console",
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"dropbox_file":                 resourceDropboxFile(),
			"dropbox_file_members":         resourceDropboxFileMembers(),
			"dropbox_folder":               resourceDropboxFolder(),
			"dropbox_paper_doc":            resourceDropboxPaperDoc(),
			"dropbox_paper_doc_users":      resourceDropboxPaperDocUsers(),
			"dropbox_paper_sharing_policy": resourceDropboxPaperSharingPolicy(),
		},

		DataSourcesMap: map[string]*schema.Resource{
			"dropbox_user_current":         dataSourceDropboxUserCurrent(),
			"dropbox_space_usage":          dataSourceDropboxSpaceUsage(),
			"dropbox_paper_folder":         dataSourceDropboxPaperFolder(),
			"dropbox_paper_sharing_policy": dataSourceDropboxPaperSharingPolicy(),
		},

		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	cfg := &ProviderConfig{
		DropboxConfig: &db.Config{Token: d.Get("access_token").(string)},
	}

	return cfg, nil
}
