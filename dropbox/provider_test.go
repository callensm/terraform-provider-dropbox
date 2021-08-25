package dropbox

import (
	"log"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/joho/godotenv"
)

var testAccProviders map[string]*schema.Provider
var testAccProvider *schema.Provider

func init() {
	godotenv.Load()

	token := os.Getenv("ACCESS_TOKEN")
	if err := os.Setenv("DROPBOX_TOKEN", token); err != nil {
		log.Fatalf("Initialization failure: Failed to set the DROPBOX_TOKEN env. %s", err)
	}

	log.Printf("DROPBOX_TOKEN successfully set to: %s", token)

	testAccProvider = Provider()
	testAccProviders = map[string]*schema.Provider{
		"dropbox": testAccProvider,
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("Provider test failure: %s", err)
	}
}

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("DROPBOX_TOKEN"); v == "" {
		t.Fatal("Precheck Failure: Environment variable DROPBOX_TOKEN is not set.")
	}
}
