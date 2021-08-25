package dropbox

import (
	"fmt"
	"testing"

	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox/paper"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDropboxResPaperSharingPolicy(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDropboxResPaperSharingConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccPaperSharingCreated("dropbox_paper_sharing_policy.foo"),
				),
			},
		},
	})
}

func testAccPaperSharingCreated(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Paper Sharing Policy Failure: %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Paper Sharing Policy Failure: ID is not set")
		}

		client := paper.New(*testAccProvider.Meta().(*ProviderConfig).DropboxConfig)

		pid := rs.Primary.Attributes["doc_id"]
		policies, err := client.DocsSharingPolicyGet(paper.NewRefPaperDoc(pid))
		if err != nil {
			return fmt.Errorf("Paper Sharing Policy Failure: %s", err)
		}

		if policies.PublicSharingPolicy.Tag == "disabled" && policies.TeamSharingPolicy.Tag == "invite_only" {
			return nil
		}

		return fmt.Errorf("Paper Sharing Policy Failure: Retrieved policies didn't match the resource input")
	}
}

const testAccDropboxResPaperSharingConfig = `
resource "dropbox_paper_doc" "doc" {
	content_file  = "${file("../token.txt")}"
	import_format = "plain_text"
}

resource "dropbox_paper_sharing_policy" "foo" {
	doc_id        = "${dropbox_paper_doc.doc.id}"
	public_policy = "disabled"
	team_policy   = "invite_only"
}
`
