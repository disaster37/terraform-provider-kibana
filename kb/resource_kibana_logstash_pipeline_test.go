package kb

import (
	"fmt"
	"testing"

	kibana "github.com/disaster37/go-kibana-rest/v7"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/pkg/errors"
)

func TestAccKibanaLogstashPipeline(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckKibanaLogstashPipelineDestroy,
		Steps: []resource.TestStep{
			{
				Config: testKibanaLogstashPipeline,
				Check: resource.ComposeTestCheckFunc(
					testCheckKibanaLogstashPipelineExists("kibana_logstash_pipeline.test"),
				),
			},
			{
				ResourceName:            "kibana_logstash_pipeline.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

func testCheckKibanaLogstashPipelineExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No logstash pipeline ID is set")
		}

		meta := testAccProvider.Meta()

		client := meta.(*kibana.Client)
		logstashPipeline, err := client.API.KibanaLogstashPipeline.Get(rs.Primary.ID)
		if err != nil {
			return err
		}
		if logstashPipeline == nil {
			return errors.Errorf("Logstash pipeline %s not found", rs.Primary.ID)
		}

		return nil
	}
}

func testCheckKibanaLogstashPipelineDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "kibana_logstash_pipeline" {
			continue
		}

		meta := testAccProvider.Meta()

		client := meta.(*kibana.Client)
		logstashPipeline, err := client.API.KibanaLogstashPipeline.Get(rs.Primary.ID)
		if err != nil {
			return err
		}
		if logstashPipeline == nil {
			return nil
		}

		return fmt.Errorf("Logstash pipeline %q still exists", rs.Primary.ID)
	}

	return nil
}

var testKibanaLogstashPipeline = `
resource "kibana_logstash_pipeline" "test" {
  name 				= "terraform-test"
  description 		= "test"
  pipeline			= "input { stdin {} } output { stdout {} }"
  settings = {
	  "queue.type" = "persisted"
  }
}
`
