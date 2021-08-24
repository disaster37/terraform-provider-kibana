package kb

import (
	"log"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

var testAccProviders map[string]terraform.ResourceProvider
var testAccProvider *schema.Provider

func init() {

	// Init provider
	testAccProvider = Provider().(*schema.Provider)
	configureFunc := testAccProvider.ConfigureFunc
	testAccProvider.ConfigureFunc = func(d *schema.ResourceData) (interface{}, error) {
		return configureFunc(d)
	}
	testAccProviders = map[string]terraform.ResourceProvider{
		"kibana": testAccProvider,
	}

}

func TestProvider(t *testing.T) {
	if err := Provider().(*schema.Provider).InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ terraform.ResourceProvider = Provider()
}

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("KIBANA_URL"); v == "" {
		t.Fatal("KIBANA_URL must be set for acceptance tests")
	}

}
