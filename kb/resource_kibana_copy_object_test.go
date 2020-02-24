package kb

import (
	"fmt"
	"testing"

	kibana "github.com/disaster37/go-kibana-rest/v7"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

func TestAccKibanaCopyObject(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckKibanaCopyObjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: testKibanaCopyObject,
				Check: resource.ComposeTestCheckFunc(
					testCheckKibanaCopyObjectExists("kibana_copy_object.test"),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testCheckKibanaCopyObjectExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No object ID is set")
		}

		// Use static value that match the current test
		objectID := "test"
		objectType := "index-pattern"
		targetSpace := "terraform-test2"

		meta := testAccProvider.Meta()

		client := meta.(*kibana.Client)
		data, err := client.API.KibanaSavedObject.Get(objectType, objectID, targetSpace)
		if err != nil {
			return err
		}

		if len(data) == 0 {
			panic(fmt.Sprintf("%+v", data))
			return errors.Errorf("Object %s not found", rs.Primary.ID)
		}

		return nil
	}
}

func testCheckKibanaCopyObjectDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		log.Infof("RS TYPE: %s", rs.Type)
		if (rs.Type != "kibana_copy_object") && (rs.Type != "kibana_object") {
			continue
		}

		log.Infof("We never delete kibana object")
		return nil

	}

	return nil

}

var testKibanaCopyObject = `
resource kibana_object "test" {
  name 				= "terraform-test"
  data				= "${file("../fixtures/test.ndjson")}"
  deep_reference	= "true"
  export_types    	= ["index-pattern"]
}

resource kibana_user_space "test" {
  name 				= "terraform-test2"
}

resource kibana_copy_object "test" {
  name 				= "terraform-test2"
  source_space		= "default"
  target_spaces		= ["${kibana_user_space.test.name}"]
  object {
	  id   = "test"
	  type = "index-pattern"
  }

  depends_on = [kibana_object.test, kibana_user_space.test]
}
`
