package kb

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceKibanaHost(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceKibanaHost,
				Check: resource.ComposeTestCheckFunc(
					testCheckDataSourceKibanaHost("kibana_host.test"),
				),
			},
		},
	})
}

func testCheckDataSourceKibanaHost(name string) resource.TestCheckFunc {
	var url, username, password resource.TestCheckFunc

	url = resource.TestCheckResourceAttr("data."+name, "url", os.Getenv("KIBANA_URL"))
	username = resource.TestCheckResourceAttr("data."+name, "username", os.Getenv("KIBANA_USERNAME"))
	password = resource.TestCheckResourceAttr("data."+name, "password", os.Getenv("KIBANA_PASSWORD"))

	return resource.ComposeAggregateTestCheckFunc(url, username, password)
}

var testDataSourceKibanaHost = `
data "kibana_host" "test" {
}
`
