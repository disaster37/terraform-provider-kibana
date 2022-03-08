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

resource kibana_user_space "test" {
  uid 				= "terraform-test"
}

resource kibana_copy_object "test" {
  name 				= "terraform-test"
  source_space		= "default"
  target_spaces		= ["${kibana_user_space.test.uid}"]
  object {
	  id   = "test"
	  type = "index-pattern"
  }
  overwrite			= true
  create_new_copies = false

  depends_on = [kibana_object.test, kibana_user_space.test]
}