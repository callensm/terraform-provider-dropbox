package dropbox

import (
	"github.com/dropbox/dropbox-sdk-go-unofficial/dropbox/users"
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceDropboxUserCurrent() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceDropboxUserCurrentRead,

		Schema: map[string]*schema.Schema{
			"account_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"display_name": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"email": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceDropboxUserCurrentRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*ProviderConfig).DropboxConfig
	client := users.New(*config)

	account, err := client.GetCurrentAccount()
	if err != nil {
		return err
	}

	d.SetId(account.AccountId)
	d.Set("account_id", account.AccountId)
	d.Set("display_name", account.Name.DisplayName)
	d.Set("email", account.Email)
	return nil
}
