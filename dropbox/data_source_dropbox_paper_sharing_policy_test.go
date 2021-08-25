package dropbox

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDropboxPaperSharingPolicy(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDropboxPaperSharingPolicyDataConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccPaperSharingPolicyExists("data.dropbox_paper_sharing_policy.foo"),
				),
			},
		},
	})
}

func testAccPaperSharingPolicyExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Paper Sharing Policy Failure: %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Paper Sharing Policy Failure: ID is not set")
		}

		if rs.Primary.Attributes["public_policy"] == "" && rs.Primary.Attributes["team_policy"] == "" {
			return fmt.Errorf("Paper Sharing Policy Failure: Both policy types were unset")
		}

		return nil
	}
}

const testAccDropboxPaperSharingPolicyDataConfig = `
data "dropbox_paper_sharing_policy" "foo" {
	doc_id = "jXl1jhXj78S7NLyloBMCB"
}
`
