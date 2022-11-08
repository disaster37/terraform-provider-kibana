package kb

import (
	"context"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/sirupsen/logrus"
	easy "github.com/t-tomalak/logrus-easy-formatter"
)

var testAccProviders map[string]*schema.Provider
var testAccProvider *schema.Provider

func init() {

	// Init logger
	logrus.SetFormatter(&easy.Formatter{
		LogFormat: "[%lvl%] %msg%\n",
	})
	logrus.SetLevel(logrus.DebugLevel)

	// Init provider
	testAccProvider = Provider()
	configureFunc := testAccProvider.ConfigureContextFunc
	testAccProvider.ConfigureContextFunc = func(ctx context.Context, rd *schema.ResourceData) (interface{}, diag.Diagnostics) {
		return configureFunc(ctx, rd)
	}
	testAccProviders = map[string]*schema.Provider{
		"kibana": testAccProvider,
	}

}

func TestProvider(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ *schema.Provider = Provider()
}

func testAccPreCheck(t *testing.T) {

	if v := os.Getenv("KIBANA_URL"); v == "" {
		t.Fatal("KIBANA_URL must be set for acceptance tests")
		panic(os.Getenv("KIBANA_URL"))
	}

}
