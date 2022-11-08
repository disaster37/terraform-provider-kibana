# kibana_object Resource Source

This resource permit to manage saved object in Kibana.
You can see the API documentation: https://www.elastic.co/guide/en/kibana/master/spaces-api.html

***Supported Kibana version:***
  - v7
  - v8

## Example Usage

It will create new object with data stored on json file.

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

## Argument Reference

***The following arguments are supported:***
  - **name**: (required) The unique name
  - **space**: (optional) The user space where to create objects
  - **data**: (required) The data to create as JSON string
  - **export_types**: (optional) The export types used to export data. It use to compare if existing is the same as in data
  - **export_objects**: (optional) The export objects used to export data. It use to compare if existing is the same as in data
  - **deep_reference**: (optional) The export deep reference. It use to compare if existing is the same as in data


## Attribute Reference

NA