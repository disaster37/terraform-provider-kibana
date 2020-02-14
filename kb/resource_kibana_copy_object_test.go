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
		objectID := "logstash-log-*"
		objectType := "index-pattern"
		targetSpace := "test"

		meta := testAccProvider.Meta()

		client := meta.(*kibana.Client)
		data, err := client.API.KibanaSavedObject.Get(objectType, objectID, targetSpace)
		if err != nil {
			return err
		}
		if len(data) == 0 {
			return errors.Errorf("Object %s not found", rs.Primary.ID)
		}

		return nil
	}
}

func testCheckKibanaCopyObjectDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "kibana_copy_object" {
			continue
		}

		log.Debugf("We never delete kibana object")
	}

	return nil
}

var testKibanaCopyObject = `
resource kibana_object "test" {
  name 				= "terraform-test"
  data				= "${file("../fixtures/index-pattern.json")}"
  deep_reference	= "true"
  export_types    	= ["index-pattern"]
}

resource kibana_space "test" {
  name 				= "terraform-test"
}

resource kibana_copy_object "test" {
  name 				= "terraform-test"
  source_space		= "default"
  target_spaces		= ["${kibana_space.test.id}"]
  object {
	  id   = "logstash-system-*"
	  type = "${kibana_object.test.export_types[0]}"
  }
}
`
