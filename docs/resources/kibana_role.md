# kibana_role Resource Source

This resource permit to manage role in Kibana.
You can see the API documentation: https://www.elastic.co/guide/en/kibana/master/role-management-api.html

***Supported Kibana version:***
  - v7

## Example Usage

It will create `role` called `terraform-test` with some privileges on index, and on Kibana feature.

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

## Argument Reference

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

## Attribute Reference

NA