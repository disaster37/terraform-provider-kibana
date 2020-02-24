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

func TestAccKibanaObject(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckKibanaObjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: testKibanaObject,
				Check: resource.ComposeTestCheckFunc(
					testCheckKibanaObjectExists("kibana_object.test"),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testCheckKibanaObjectExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No object ID is set")
		}

		// Use static value that match the current test
		deepReference := true
		exportObject := map[string]string{}
		exportObject["id"] = "logstash-log-*"
		exportObject["type"] = "index-pattern"
		exportObjects := []map[string]string{exportObject}
		space := "default"

		meta := testAccProvider.Meta()

		client := meta.(*kibana.Client)
		data, err := client.API.KibanaSavedObject.Export(nil, exportObjects, deepReference, space)
		if err != nil {
			return err
		}
		if len(data) == 0 {
			return errors.Errorf("Object %s not found", rs.Primary.ID)
		}

		return nil
	}
}

func testCheckKibanaObjectDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "kibana_object" {
			continue
		}

		log.Debugf("We never delete kibana object")
	}

	return nil
}

var testKibanaObject = `
resource "kibana_object" "test" {
  name 				= "terraform-test"
  data				= "${file("../fixtures/index-pattern.json")}"
  deep_reference	= "true"
  export_types    	= ["index-pattern"]
}
`
