package dropbox

import (
	"fmt"
	"testing"

	db "github.com/dropbox/dropbox-sdk-go-unofficial/dropbox"
	"github.com/dropbox/dropbox-sdk-go-unofficial/dropbox/paper"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccDropboxPaperDoc_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccDropboxPaperDocConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccPaperDocCreated("dropbox_paper_doc.doc"),
				),
			},
		},
	})
}

func TestAccDropboxPaperDoc_foldered(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccDropboxPaperDocFolderConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccPaperDocCreated("dropbox_paper_doc.foldered_doc"),
				),
			},
		},
	})
}

func testAccPaperDocCreated(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Paper Doc Failure: %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Paper Doc Failure: ID is not set")
		}

		config := testAccProvider.Meta().(*ProviderConfig).DropboxConfig
		client := paper.New(*config)

		opts := &paper.ListPaperDocsArgs{
			Limit:     1000,
			FilterBy:  &paper.ListPaperDocsFilterBy{Tagged: db.Tagged{Tag: "docs_created"}},
			SortBy:    &paper.ListPaperDocsSortBy{Tagged: db.Tagged{Tag: "created"}},
			SortOrder: &paper.ListPaperDocsSortOrder{Tagged: db.Tagged{Tag: "descending"}},
		}
		results, err := client.DocsList(opts)
		if err != nil {
			return fmt.Errorf("Paper Doc Failure: %s", err)
		}

		for _, id := range results.DocIds {
			if id == rs.Primary.Attributes["doc_id"] {
				return nil
			}
		}

		return fmt.Errorf("Paper Doc Failure: Document with ID %s was not found in fetched list", rs.Primary.Attributes["doc_id"])
	}
}

const testAccDropboxPaperDocConfig = `
resource "dropbox_paper_doc" "doc" {
	content_file  = "${file("../token.txt")}"
	import_format = "plain_text"
}
`

const testAccDropboxPaperDocFolderConfig = `
resource "dropbox_paper_doc" "foldered_doc" {
	content_file  = "${file("../token.txt")}"
	parent_folder = "e.1gg8YzoPEhbTkrhvQwJ2zz3QnCdRv14CYjZI6kJODQUKnl1Usxt7"
	import_format = "plain_text"
}
`
