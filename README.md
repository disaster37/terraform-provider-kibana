# terraform-provider-kibana

[![CircleCI](https://circleci.com/gh/disaster37/terraform-provider-kibana/tree/7.x.svg?style=svg)](https://circleci.com/gh/disaster37/terraform-provider-kibana/tree/7.x)
[![Go Report Card](https://goreportcard.com/badge/github.com/disaster37/terraform-provider-kibana)](https://goreportcard.com/report/github.com/disaster37/terraform-provider-kibana)
[![GoDoc](https://godoc.org/github.com/disaster37/terraform-provider-kibana?status.svg)](http://godoc.org/github.com/disaster37/terraform-provider-kibana)
[![codecov](https://codecov.io/gh/disaster37/terraform-provider-kibana/branch/7.x/graph/badge.svg)](https://codecov.io/gh/disaster37/terraform-provider-kibana/branch/7.x)

This is a terraform provider that lets you provision kibana resources, compatible with v7 of kibana.
For Kibana 7, you need to use branch and release 7.x

## Installation

[Download a binary](https://github.com/disaster37/terraform-provider-kibana/releases), and put it in a good spot on your system. Then update your `~/.terraformrc` to refer to the binary:

```hcl
providers {
  kibana = "/path/to/terraform-provider-kibana"
}
```

See [the docs for more information](https://www.terraform.io/docs/plugins/basics.html).

## Usage

### Provider

The Kibana provider is used to interact with the
resources supported by Elasticsearch. The provider needs
to be configured with an endpoint URL before it can be used.

***Sample:***
```tf
provider "kibana" {
    url     = "http://kibana.company.com:5601"
    username = "elastic"
    password = "changeme"
}
```

***The following arguments are supported:***
- **url**: (required) The endpoint Kibana URL.
- **username**: (optional) The username to connect on it.
- **password**: (optional) The password to connect on it.
- **insecure**: (optional) To disable the certificate check.
- **cacert_files**: (optional) The list of CA contend to use if you use custom PKI.
- **retry**: (optional) The number of time you should to retry connexion befaore exist with error. Default to `6`.
- **wait_before_retry**: (optional) The number of time in second we wait before each connexion retry. Default to `10`.

___


### Role resource

This resource permit to manage role in Kibana.
You can see the API documentation: https://www.elastic.co/guide/en/kibana/master/role-management-api.html

***Supported Kibana version:***
  - v7

***Sample:***
```tf
resource kibana_role "test" {
  name 				= "terraform-test"
  elasticsearch {
	indices {
		names 		= ["logstash-*"]
		privileges 	= ["read"]
	}
	indices {
		names 		= ["logstash-*"]
		privileges 	= ["read2"]
	}
	cluster = ["all"]
  }
  kibana {
	  features {
		  name 			= "dashboard"
		  permissions 	= ["read"]
	  }
	  features {
		  name 			= "discover"
		  permissions 	= ["read"]
	  }
	  spaces = ["default"]
  }
}
```

***The following arguments are supported:***
  - **name**: (required) The role name to create
  - **elasticsearch**: (optional) The elasticsearch permission object
  - **kibana**: (optional) The kibana permission object
  - **metadata**: (optional) A string as JSON object meta-data. Within the metadata object, keys that begin with _ are reserved for system usage.

***Elasticsearch permission object***:
  - **cluster**: (optional) A list of cluster privileges. These privileges define the cluster level actions that users with this role are able to execute.
  - **run_as**: (optional) A list of users that the owners of this role can impersonate.
  - **indices**: (optional) A list of indices permissions entries. Look the indice object below.

***Kibana permission object***:
  - **base**: (optional) A base privilege. When specified, the base must be ["all"] or ["read"]. When the base privilege is specified, you are unable to use the feature section. "all" grants read/write access to all Kibana features for the specified spaces. "read" grants read-only access to all Kibana features for the specified spaces.
  - **spaces**: (required) The spaces to apply the privileges to. To grant access to all spaces, set to ["*"]
  - **features**: (optional) Contains privileges for specific features. When the feature privileges are specified, you are unable to use the base section

***Indice object***:
  - **names**: (required) A list of indices (or index name patterns) to which the permissions in this entry apply.
  - **privileges**: (required) A list of The index level privileges that the owners of the role have on the specified indices.
  - **query**: (optional) A search query that defines the documents the owners of the role have read access to. A document within the specified indices must match this query in order for it to be accessible by the owners of the role. It's a string or a string as JSON object.
  - **field_security**: (optional) The document fields that the owners of the role have read access to. It's a string as JSON object

___

### User space management

This resource permit to manage user space in Kibana.
You can see the API documentation: https://www.elastic.co/guide/en/kibana/master/spaces-api.html

***Supported Kibana version:***
  - v7

***Sample:***
```tf
resource kibana_user_space "test" {
  name 				= "terraform-test"
  description 		= "test"
  initials			= "tt"
  color				= "#000000"
  disabled_features = ["canvas", "maps", "advancedSettings", "indexPatterns", "graph", "monitoring", "ml", "apm", "infrastructure", "logs", "siem"]
}
```

***The following arguments are supported:***
  - **name**: (required) The user space name to create
  - **description**: (optional) The description for user space
  - **disabled_features**: (optional) The list of features you should disabled for this user space.
  - **initials**: (optional) The initial for user space
  - **color**: (optional) The color for user space


### Saved object management

This resource permit to manage saved object in Kibana.
You can see the API documentation: https://www.elastic.co/guide/en/kibana/master/spaces-api.html

***Supported Kibana version:***
  - v7

***Sample:***
```tf
resource kibana_object "test" {
  name 				= "terraform-test"
  data				= "${file("../fixtures/index-pattern.json")}"
  deep_reference	= "true"
  export_objects {
	  id = "logstash-log-*"
	  type = "index-pattern"
  }
}
```

***The following arguments are supported:***
  - **name**: (required) The unique name
  - **space**: (optional) The user space where to create objects
  - **data**: (required) The data to create as JSON string
  - **export_types**: (optional) The export types used to export data. It use to compare if existing is the same as in data
  - **export_objects**: (optional) The export objects used to export data. It use to compare if existing is the same as in data
  - **deep_reference**: (optional) The export deep reference. It use to compare if existing is the same as in data

---

### Copy saved object

This resource permit to copy objects from space to another spaces.
You can see the API documentation: https://www.elastic.co/guide/en/kibana/master/spaces-api.html

***Supported Kibana version:***
  - v7

***Sample:***
```tf
resource kibana_copy_object "test" {
  name 				= "terraform-test"
  source_space		= "default"
  target_spaces		= ["Team A"]
  object {
	  id   = "logstash-system-*"
	  type = "index-pattern"
  }
}
```

***The following arguments are supported:***
  - **name**: (required) The unique name
  - **source_space**: (optional) The user space from copy objects. Default to `default`
  - **target_spaces**: (required) The list of space where to copy objects
  - **overwrite**: (optional) Overwrite existing objects. Default to `true`
  - **object**: (optional) The list of object you should to copy
  - **include_reference**: (optional) Include reference when copy objects. Default to `true`
  - **force_update**: (optional) Force to copy objects each time you apply. Default to `true`

***object:***
  - **id**: (required) The object ID
  - **type**: (required) The object type

---

### Logstash pipeline management

This resource permit to manage logstash pipeline in Kibana.
You can see the API documentation: https://www.elastic.co/guide/en/kibana/master/logstash-configuration-management-api.html

***Supported Kibana version:***
  - v7

***Sample:***
```tf
resource kibana_logstash_pipeline "test" {
  name 				= "terraform-test"
  description 		= "test"
  pipeline			= "input { stdin {} } output { stdout {} }"
  settings = {
	  "queue.type" = "persisted"
  }
}
```

***The following arguments are supported:***
  - **name**: (required) The unique name of logstash pipeline
  - **description**: (optional) The logstash pipeline description
  - **pipeline**: (required) The pipeline specification as JSON string.
  - **settings**: (optional) The extra logstash pipeline settings, as map of string.

***Computed field***
  - **username**: The username that create the logstash pipeline

---

## Development

### Requirements

* [Golang](https://golang.org/dl/) >= 1.11
* [Terrafrom](https://www.terraform.io/) >= 0.12


```
go build -o /path/to/binary/terraform-provider-kibana
```

## Licence

See LICENSE.

## Contributing

1. Fork it ( https://github.com/disaster37/terraform-provider-kibana/fork )
2. Go to the right branch (7.x for Kibana 7) (`git checkout 7.x`)
3. Create your feature branch (`git checkout -b my-new-feature`)
4. Add feature, add acceptance test and tets your code (`KIBANA_URL=http://127.0.0.1:5601 KIBANA_USERNAME=elastic KIBANA_PASSWORD=changeme make testacc`)
5. Commit your changes (`git commit -am 'Add some feature'`)
6. Push to the branch (`git push origin my-new-feature`)
7. Create a new Pull Request