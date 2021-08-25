package dropbox

import (
	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox/users"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceDropboxUserCurrent() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceDropboxUserCurrentRead,

		Schema: map[string]*schema.Schema{
			"account_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"account_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"display_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"email": {
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
	d.Set("account_type", account.AccountType.Tag)
	d.Set("display_name", account.Name.DisplayName)
	d.Set("email", account.Email)
	return nil
}
