terraform {
  required_providers {
    kibana = {
      source = "disaster37/kibana"
    }
  }
}

provider "kibana" {
    url      = "http://kibana:5601"
    username = "elastic"
    password = "changeme"
}


resource "kibana_object" "test" {
  name 				= "terraform-test"
  data				= "{\"id\": \"test\", \"type\": \"index-pattern\",\"attributes\": {\"title\": \"test\"}}"
  deep_reference	= "true"
  export_types    	= ["index-pattern"]
}