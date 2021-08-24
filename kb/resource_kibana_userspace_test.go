package kb

import (
	"fmt"
	"testing"

	kibana "github.com/disaster37/go-kibana-rest/v7"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/pkg/errors"
)

func TestAccKibanaUserSpace(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckKibanaUserSpaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testKibanaUserSpace,
				Check: resource.ComposeTestCheckFunc(
					testCheckKibanaUserSpaceExists("kibana_user_space.test"),
				),
			},
			{
				ResourceName:            "kibana_user_space.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

func testCheckKibanaUserSpaceExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No user space ID is set")
		}

		meta := testAccProvider.Meta()

		client := meta.(*kibana.Client)
		userSpace, err := client.API.KibanaSpaces.Get(rs.Primary.ID)
		if err != nil {
			return err
		}
		if userSpace == nil {
			return errors.Errorf("User space %s not found", rs.Primary.ID)
		}

		return nil
	}
}

func testCheckKibanaUserSpaceDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "kibana_user_space" {
			continue
		}

		meta := testAccProvider.Meta()

		client := meta.(*kibana.Client)
		userSpace, err := client.API.KibanaSpaces.Get(rs.Primary.ID)
		if err != nil {
			return err
		}
		if userSpace == nil {
			return nil
		}

		return fmt.Errorf("User space %q still exists", rs.Primary.ID)
	}

	return nil
}

var testKibanaUserSpace = `
resource "kibana_user_space" "test" {
  name 				= "terraform-test"
  description 		= "test"
  initials			= "tt"
  color				= "#000000"
  disabled_features = ["canvas", "maps", "advancedSettings", "indexPatterns", "graph", "monitoring", "ml", "apm", "infrastructure", "logs", "siem"]
}
`
