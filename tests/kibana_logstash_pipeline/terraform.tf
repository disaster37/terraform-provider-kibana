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


resource kibana_logstash_pipeline "test" {
    name = "terraform-test"
    description = "test"
    pipeline = "input { stdin{} } output { stdout{} }"
    settings = {
    "queue.type" = "persisted"
    }
}