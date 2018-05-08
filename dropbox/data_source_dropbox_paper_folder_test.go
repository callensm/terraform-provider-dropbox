package dropbox

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccDropboxPaperFolder(t *testing.T) {
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

// TODO: Insert real document ID
const testAccDropboxFolderConfig = `
	data "dropbox_paper_folder" "foo" {
		doc_id = "abc1234"
	}
`
