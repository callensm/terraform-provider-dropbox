package dropbox

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccDropboxUserCurrent(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccDropboxUserCurrentDataConfig,
				Check:  resource.ComposeTestCheckFunc(),
			},
		},
	})
}

const testAccDropboxUserCurrentDataConfig = `
data "dropbox_user_current" "user" {}
`
