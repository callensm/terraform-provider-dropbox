package dropbox

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDropboxPaperDocUsers(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDropboxPaperUsersConfig,
				Check:  resource.ComposeTestCheckFunc(),
			},
		},
	})
}

const testAccDropboxPaperUsersConfig = `
resource "dropbox_paper_doc_users" "foo" {
	doc_id  = "oeWK68vUIXUnDe3r5H5wo"
	quiet   = false
	members = [
		{
			identity    = "callensmatt@gmail.com"
			permissions = "edit"
		}
	]
}
`
