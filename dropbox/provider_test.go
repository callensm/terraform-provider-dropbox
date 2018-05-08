package dropbox

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

var testAccProviders map[string]terraform.ResourceProvider
var testAccProvider *schema.Provider

func init() {
	testAccProvider = Provider().(*schema.Provider)
	testAccProviders = map[string]terraform.ResourceProvider{
		"dropbox": testAccProvider,
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().(*schema.Provider).InternalValidate(); err != nil {
		t.Fatalf("Provider test failure: %s", err)
	}
}

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("DROPBOX_TOKEN"); v == "" {
		t.Log("Precheck warning: Environment variable DROPBOX_TOKEN is not set. Setting now based on token.txt for testing only.")

		token, err := ioutil.ReadFile("../token.txt")
		if err != nil {
			t.Fatalf("Precheck failure: Couldn't read access token file. %s", err)
		}

		err = os.Setenv("DROPBOX_TOKEN", string(token))
		if err != nil {
			t.Fatalf("Precheck failure: Failed to set the DROPBOX_TOKEN env. %s", err)
		}
	}
}
