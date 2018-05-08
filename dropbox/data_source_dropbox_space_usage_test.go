package dropbox

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccDropboxSpaceUsage(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccDropboxSpaceUsageDataConfig,
				Check:  resource.ComposeTestCheckFunc(),
			},
		},
	})
}

const testAccDropboxSpaceUsageDataConfig = `
data "dropbox_space_usage" "usage" {}
`
