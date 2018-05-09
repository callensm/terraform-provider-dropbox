package dropbox

import (
	"fmt"
	"testing"

	"github.com/dropbox/dropbox-sdk-go-unofficial/dropbox/files"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccDropboxFolder_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccDropboxFolderConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccFolderCreated("dropbox_folder.foo"),
				),
			},
		},
	})
}

func TestAccDropboxFolder_nested(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccDropboxFolderNestedConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccFolderCreated("dropbox_folder.test_B"),
				),
			},
		},
	})
}

func testAccFolderCreated(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Folder Failure: %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Folder Failure: ID is not set")
		}

		config := testAccProvider.Meta().(*ProviderConfig).DropboxConfig
		client := files.New(*config)
		path := rs.Primary.Attributes["path"]

		opts := &files.ListFolderArg{
			Path:      path,
			Recursive: true,
		}
		results, err := client.ListFolder(opts)
		if err != nil {
			return fmt.Errorf("Folder Failure: Issue listing folders %s", err)
		}

		for _, e := range results.Entries {
			if e.(*files.FolderMetadata).PathDisplay == path {
				return nil
			}
		}

		return fmt.Errorf("Folder Failure: No folder system was found that matched %s", path)
	}
}

const testAccDropboxFolderConfig = `
resource "dropbox_folder" "foo" {
	path        = "/terraform-created"
	auto_rename = false
}
`

const testAccDropboxFolderNestedConfig = `
resource "dropbox_folder" "test_A" {
	path = "/terraform-test-A"
}

resource "dropbox_folder" "test_B" {
	path = "${dropbox_folder.test_A.path}/test-B"
}
`
