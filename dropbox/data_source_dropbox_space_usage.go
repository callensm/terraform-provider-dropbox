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
			"allocation": &schema.Schema{
				Type:     schema.TypeMap,
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

	var allocated uint64
	isTeam := false
	if indiv := usage.Allocation.Individual; indiv != nil {
		allocated = indiv.Allocated
	} else {
		isTeam = true
		allocated = usage.Allocation.Team.Allocated
	}

	d.SetId(fmt.Sprintf("%d:%v", usage.Used, usage.Allocation))
	d.Set("used", usage.Used)
	d.Set("allocation", allocated)
	d.Set("is_team_allocation", isTeam)
	return nil
}
