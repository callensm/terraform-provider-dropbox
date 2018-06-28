package dropbox

import (
	"fmt"
	"testing"

	"github.com/dropbox/dropbox-sdk-go-unofficial/dropbox/sharing"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccDropboxFolderMember(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccDropboxFolderMemberConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccFolderMembersSet("dropbox_folder_members.foo"),
				),
			},
		},
	})
}

func testAccFolderMembersSet(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Folder Member Failure: %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Folder Member Failure: ID is not set")
		}

		config := testAccProvider.Meta().(*ProviderConfig).DropboxConfig
		client := sharing.New(*config)
		memberCount := len(rs.Primary.Attributes["members"])

		results, err := client.ListFolderMembers(&sharing.ListFolderMembersArgs{
			ListFolderMembersCursorArg: *sharing.NewListFolderMembersCursorArg(),
			SharedFolderId:             rs.Primary.Attributes["folder_id"],
		})

		if err != nil {
			return fmt.Errorf("Folder Member Failure: Issue listing folders %s", err)
		}

		if len(results.Users) != memberCount {
			return fmt.Errorf("Folder Member Failure: Retrieve member count doesn't match argued count")
		}

		return nil
	}
}

const testAccDropboxFolderMemberConfig = `
resource "dropbox_folder" "src" {
	path        = "/terraform-member-test"
	auto_rename = false
}

resource "dropbox_folder_members" "foo" {
	folder_id = "${dropbox_folder.src.folder_id}"
	members   = [
		{
			email        = "test@example.com"
			access_level = "viewer"
		},
		{
			account_id = "agqorghqovn34"
		}
	]
	quiet     = true
}
`
