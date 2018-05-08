package dropbox

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccDropboxUserCurrent(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccDropboxUserCurrentDataConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccUserCurrentExists("data.dropbox_user_current.user"),
				),
			},
		},
	})
}

func testAccUserCurrentExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("User Current Failure: %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("User Current Failure: ID is not set")
		}

		for attr, val := range rs.Primary.Attributes {
			if val == "" {
				return fmt.Errorf("User Current Failure: Attribute %s was not found", attr)
			}
		}

		fmt.Printf("User: %+v\n", rs.Primary.Attributes)
		return nil
	}
}

const testAccDropboxUserCurrentDataConfig = `
data "dropbox_user_current" "user" {}
`
