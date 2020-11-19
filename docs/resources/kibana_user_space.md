# kibana_user_space Resource Source

This resource permit to manage user space in Kibana.
You can see the API documentation: https://www.elastic.co/guide/en/kibana/master/spaces-api.html

***Supported Kibana version:***
  - v7

## Example Usage

It will create `space` called `terraform-test` with some features disabled.

```tf
resource kibana_user_space "test" {
  name 				= "terraform-test"
  description 		= "test"
  initials			= "tt"
  color				= "#000000"
  disabled_features = ["canvas", "maps", "advancedSettings", "indexPatterns", "graph", "monitoring", "ml", "apm", "infrastructure", "logs", "siem"]
}
```

## Argument Reference

***The following arguments are supported:***
  - **name**: (required) The user space name to create
  - **description**: (optional) The description for user space
  - **disabled_features**: (optional) The list of features you should disabled for this user space.
  - **initials**: (optional) The initial for user space
  - **color**: (optional) The color for user space

## Attribute Reference

NA