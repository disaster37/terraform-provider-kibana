package kb

import (
	"fmt"
	"testing"

	kibana7 "github.com/disaster37/go-kibana-rest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
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

		switch meta.(type) {
		// v7
		case *kibana7.Client:
			client := meta.(*kibana7.Client)
			role, err := client.API.KibanaRoleManagement.Get(rs.Primary.ID)
			if err != nil {
				return err
			}
			if role == nil {
				return errors.Errorf("role %s not found", rs.Primary.ID)
			}
		default:
			return errors.New("Role is only supported by the kibana library >= v6!")
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

		switch meta.(type) {
		// v7
		case *kibana7.Client:
			client := meta.(*kibana7.Client)
			role, err := client.API.KibanaRoleManagement.Get(rs.Primary.ID)
			if err != nil {
				return err
			}
			if role == nil {
				return nil
			}
		default:
			return errors.New("Role is only supported by the kibana library >= v6!")
		}

		return fmt.Errorf("Role %q still exists", rs.Primary.ID)
	}

	return nil
}

var testKibanaRole = `
resource "kibana_role" "test" {
  name 				= "terraform-test"
  elasticsearch {
	indices {
		names 		= ["logstash-*"]
		privileges 	= ["read"]
	}
	indices {
		names 		= ["logstash-*"]
		privileges 	= ["read2"]
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
