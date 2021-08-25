package dropbox

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDropboxPaperFolder(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDropboxPaperFolderDataConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccPaperFolderExists("data.dropbox_paper_folder.foo"),
				),
			},
		},
	})
}

func testAccPaperFolderExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Paper Folder Failure: %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Paper Folder Failure: ID is not set")
		}

		if num, _ := strconv.Atoi(rs.Primary.Attributes["folders.#"]); num < 1 {
			return fmt.Errorf("Paper Folder Failure: Should find at least 1 folder but instead found %d", num)
		}

		return nil
	}
}

const testAccDropboxPaperFolderDataConfig = `
data "dropbox_paper_folder" "foo" {
	doc_id = "jXl1jhXj78S7NLyloBMCB"
}
`
