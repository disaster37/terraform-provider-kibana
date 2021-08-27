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


resource "kibana_role" "test" {
  name 				= "terraform-test"
  kibana {
      base   = ["read"]
      spaces = ["default"]
  }
}