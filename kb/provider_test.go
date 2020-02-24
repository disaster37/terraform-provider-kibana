package kb

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
)

var testAccProviders map[string]terraform.ResourceProvider
var testAccProvider *schema.Provider

func init() {

	// Init logger
	logrus.SetFormatter(new(prefixed.TextFormatter))
	logrus.SetLevel(logrus.InfoLevel)

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
