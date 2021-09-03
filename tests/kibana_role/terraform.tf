terraform {
  required_providers {
    kibana = {
      source = "disaster37/kibana"
      version = "1.0.0"
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

resource kibana_role "my_kibana_role" {
  name = "my_kibana_role"
  elasticsearch {
    indices {
      names = ["*"]
      privileges = ["read"]
      field_security = jsonencode({"grant": ["my_field"]})
    }
  }
}