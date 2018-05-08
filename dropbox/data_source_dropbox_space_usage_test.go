package dropbox

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccDropboxSpaceUsage(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccDropboxSpaceUsageDataConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccValidUsage("data.dropbox_space_usage.usage"),
				),
			},
		},
	})
}

func testAccValidUsage(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Space Usage Failure: %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Space Usage Failure: ID is not set")
		}

		if v, _ := strconv.Atoi(rs.Primary.Attributes["used"]); v == 0 {
			return fmt.Errorf("Space Usage Failure: 0B or invalid usage amount")
		}

		fmt.Printf("Current Usage: %v\n", rs.Primary.Attributes["used"])
		return nil
	}
}

const testAccDropboxSpaceUsageDataConfig = `
data "dropbox_space_usage" "usage" {}
`
