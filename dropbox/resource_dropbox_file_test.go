package dropbox

import (
	"fmt"
	"testing"

	"github.com/dropbox/dropbox-sdk-go-unofficial/dropbox/files"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccDropboxFile_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccDropboxFileConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccFileCreated("dropbox_file.foo"),
				),
			},
		},
	})
}

func TestAccDropboxFile_folder(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccDropboxFileFolderConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccFileCreated("dropbox_file.bar"),
				),
			},
		},
	})
}

func testAccFileCreated(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("File Failure: %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("File Failure: ID is not set")
		}

		config := testAccProvider.Meta().(*ProviderConfig).DropboxConfig
		client := files.New(*config)
		path := rs.Primary.Attributes["path"]

		opts := files.NewGetMetadataArg(path)
		result, err := client.GetMetadata(opts)
		if err != nil {
			return fmt.Errorf("File Failure: %s", err)
		}

		if result.(*files.FileMetadata).ContentHash != rs.Primary.Attributes["hash"] {
			return fmt.Errorf("File Failure: Content hashes for created and found don't match")
		}

		return nil
	}
}

const testAccDropboxFileConfig = `
resource "dropbox_file" "foo" {
	content = "${file("../Makefile")}"
	path    = "/Makefile"
	mute    = true
}
`

const testAccDropboxFileFolderConfig = `
resource "dropbox_folder" "dest" {
	path        = "/test"
	auto_rename = false
}

resource "dropbox_file" "bar" {
	content = "${file("../Makefile")}"
	path    = "/test/Makefile"
	mode    = "add"
	mute    = true
}
`
