# kibana_host Data Source

This resource permit to retrieve the Kibana connection settings.

***Supported Kibana version:***

- v7

## Example Usage

```tf
data kibana_host "test" {
}
```

## Argument Reference

NA

## Attribute Reference

- **url**: The Kibana URL
- **username**: The username to use to connect to Kibana. If empty, no authentication is needed
- **password**: The password to use to connect to Kibana. If empty, no authentication is needed
