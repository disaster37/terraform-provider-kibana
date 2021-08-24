package kb

import (
	"fmt"
	"testing"

	kibana "github.com/disaster37/go-kibana-rest/v7"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/pkg/errors"
)

func TestAccKibanaRole(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckKibanaRoleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testKibanaRole,
				Check: resource.ComposeTestCheckFunc(
					testCheckKibanaRoleExists("kibana_role.test"),
				),
			},
			{
				ResourceName:            "kibana_role.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"elasticsearch", "kibana", "metadata"},
			},
		},
	})
}

func testCheckKibanaRoleExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No role ID is set")
		}

		meta := testAccProvider.Meta()

		client := meta.(*kibana.Client)
		role, err := client.API.KibanaRoleManagement.Get(rs.Primary.ID)
		if err != nil {
			return err
		}
		if role == nil {
			return errors.Errorf("role %s not found", rs.Primary.ID)
		}

		return nil
	}
}

func testCheckKibanaRoleDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "kibana_role" {
			continue
		}

		meta := testAccProvider.Meta()

		client := meta.(*kibana.Client)
		role, err := client.API.KibanaRoleManagement.Get(rs.Primary.ID)
		if err != nil {
			return err
		}
		if role == nil {
			return nil
		}

		return fmt.Errorf("Role %q still exists", rs.Primary.ID)
	}

	return nil
}

var testKibanaRole = `
resource kibana_role "test" {
  name 				= "terraform-test"
  elasticsearch {
	indices {
		names 		= ["logstash-*"]
		privileges 	= ["read"]
	}
	indices {
		names 		= ["logstash-*"]
		privileges 	= ["read"]
	}
	cluster = ["all"]
  }
  kibana {
	  features {
		  name 			= "dashboard"
		  permissions 	= ["read"]
	  }
	  features {
		  name 			= "discover"
		  permissions 	= ["read"]
	  }
	  spaces = ["default"]
  }
}
`
