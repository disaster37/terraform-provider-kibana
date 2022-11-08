# kibana_copy_object Resource Source

This resource permit to copy objects from space to another spaces.
You can see the API documentation: https://www.elastic.co/guide/en/kibana/master/spaces-api.html

***Supported Kibana version:***
  - v7
  - v8

## Example Usage

It will copy `logstash-system-*` index-pattern from `default` space to  `Team_A` space.

```tf
resource kibana_copy_object "test" {
  name 				= "terraform-test"
  source_space		= "default"
  target_spaces		= ["Team_A"]
  object {
	  id   = "logstash-system-*"
	  type = "index-pattern"
  }
}
```

## Argument Reference

***The following arguments are supported:***
  - **name**: (required) The unique name
  - **source_space**: (optional) The user space from copy objects. Default to `default`
  - **target_spaces**: (required) The list of space where to copy objects
  - **overwrite**: (optional) Overwrite existing objects. Default to `false`
  - **create_new_copies**: (optional)  Creates new copies of saved objects, regenerates each object ID, and resets the origin. Default to `true`.
  - **object**: (optional) The list of object you should to copy
  - **include_reference**: (optional) Include reference when copy objects. Default to `true`
  - **force_update**: (optional) Force to copy objects each time you apply. Default to `true`

***object:***
  - **id**: (required) The object ID
  - **type**: (required) The object type

## Attribute Reference

NA