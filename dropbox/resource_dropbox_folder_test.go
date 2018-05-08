package dropbox

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccDropboxFolder(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccDropboxFolderConfig,
				Check:  resource.ComposeTestCheckFunc(),
			},
		},
	})
}

// TODO: Insert valid path variables for config
const testAccDropboxFolderConfig = `
resource "dropbox_folder" "foo" {
	path        = ""
	auto_rename = false
}
`
