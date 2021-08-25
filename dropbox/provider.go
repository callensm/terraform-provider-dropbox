package dropbox

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Provider for the Dropbox API in Terraform
// returns a terraform.ResourceProvider
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"access_token": {
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
			"dropbox_folder_members":       resourceDropboxFolderMembers(),
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

		ConfigureFunc: providerConfig,
	}
}
