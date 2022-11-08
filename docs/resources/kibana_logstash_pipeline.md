# kibana_logstash_pipeline Resource Source

This resource permit to manage logstash pipeline in Kibana.
You can see the API documentation: https://www.elastic.co/guide/en/kibana/master/logstash-configuration-management-api.html

***Supported Kibana version:***
  - v7
  - v8

## Example Usage

It will create `pipeline` called `terraform-test` with some logstash rules.

```tf
resource kibana_logstash_pipeline "test" {
  name 				= "terraform-test"
  description 		= "test"
  pipeline			= "input { stdin {} } output { stdout {} }"
  settings {
	  queue_type = "persisted"
    pipeline_workers = 2
  }
}
```

## Argument Reference

***The following arguments are supported:***
  - **name**: (required) The unique name of logstash pipeline
  - **description**: (optional) The logstash pipeline description
  - **pipeline**: (required) The pipeline specification as JSON string.
  - **settings**: (optional) The extra logstash pipeline settings, as object.

*** Settings object:***
  - **pipeline_workers**: (optional)
  - **pipeline_batch_size**: (optional)
  - **pipeline_batch_delay**: (optional)
  - **pipeline_ecs_compatibility**: (optional)
  - **pipeline_ordored**: (optional)
  - **queue_type**: (optional)
  - **queue_max_bytes**: (optional)
  - **queue_checkpoint_writes**: (optional)


## Attribute Reference

***Computed field***
  - **username**: The username that create the logstash pipeline