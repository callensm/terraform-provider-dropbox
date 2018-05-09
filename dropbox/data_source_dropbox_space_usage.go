package dropbox

import (
	"fmt"

	"github.com/dropbox/dropbox-sdk-go-unofficial/dropbox/users"
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceDropboxSpaceUsage() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceDropboxSpaceUsageRead,

		Schema: map[string]*schema.Schema{
			"used": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
			"allocated": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
			"is_team_allocation": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func dataSourceDropboxSpaceUsageRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*ProviderConfig).DropboxConfig
	client := users.New(*config)

	usage, err := client.GetSpaceUsage()
	if err != nil {
		return err
	}

	used := usage.Used
	isTeam := false

	var allocated uint64
	if tag := usage.Allocation.Tag; tag == "individual" {
		allocated = usage.Allocation.Individual.Allocated
	} else {
		isTeam = true
		used = usage.Allocation.Team.Used
		allocated = usage.Allocation.Team.Allocated
	}

	d.SetId(fmt.Sprintf("%d:%d", used, allocated))
	d.Set("used", usage.Used)
	d.Set("allocated", allocated)
	d.Set("is_team_allocation", isTeam)
	return nil
}
