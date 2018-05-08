package dropbox

import (
	"io/ioutil"
	"log"
	"os"
	"strings"
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

	token, err := ioutil.ReadFile("../token.txt")
	if err != nil {
		log.Fatalf("Initialization failure: Couldn't read access token file. %s", err)
	}

	err = os.Setenv("DROPBOX_TOKEN", strings.Replace(string(token), "\n", "", 1))
	if err != nil {
		log.Fatalf("Initialization failure: Failed to set the DROPBOX_TOKEN env. %s", err)
	}

	log.Printf("DROPBOX_TOKEN successfully set to: %s", token)
}

func TestProvider(t *testing.T) {
	if err := Provider().(*schema.Provider).InternalValidate(); err != nil {
		t.Fatalf("Provider test failure: %s", err)
	}
}

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("DROPBOX_TOKEN"); v == "" {
		t.Fatal("Precheck Failure: Environment variable DROPBOX_TOKEN is not set.")
	}
}
