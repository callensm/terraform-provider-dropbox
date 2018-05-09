package dropbox

import (
	"fmt"

	"github.com/dropbox/dropbox-sdk-go-unofficial/dropbox/paper"
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceDropboxPaperSharingPolicy() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceDropboxPaperSharingPolicyRead,

		Schema: map[string]*schema.Schema{
			"doc_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"public_policy": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"team_policy": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func dataSourceDropboxPaperSharingPolicyRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*ProviderConfig).DropboxConfig
	client := paper.New(*config)

	opts := paper.NewRefPaperDoc(d.Get("doc_id").(string))
	policy, err := client.DocsSharingPolicyGet(opts)
	if err != nil {
		return fmt.Errorf("Sharing Policy Data Failure: %s", err)
	}

	d.SetId(fmt.Sprintf("policy_ds:%s", d.Get("doc_id").(string)))

	if public := policy.PublicSharingPolicy; public != nil {
		d.Set("public_policy", public.Tag)
	}

	if team := policy.TeamSharingPolicy; team != nil {
		d.Set("team_policy", team.Tag)
	}

	return nil
}
