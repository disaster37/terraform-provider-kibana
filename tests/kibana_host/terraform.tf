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


data "kibana_host" "test" {
}


output "url" {
  value = data.kibana_host.test.url
}

output "username" {
  value = data.kibana_host.test.username
}

output "password" {
  value = data.kibana_host.test.password
}