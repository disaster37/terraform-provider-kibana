# Kibana Provider

This is a terraform provider that lets you provision kibana resources, compatible with v7, v8 of kibana.

## Example Usage

The Kibana provider is used to interact with the
resources supported by Kibana. The provider needs
to be configured with an endpoint URL before it can be used.

***Sample:***
```tf
provider "kibana" {
    url     = "http://kibana.company.com:5601"
    username = "elastic"
    password = "changeme"
}
```

## Argument Reference

***The following arguments are supported:***
- **url**: (required) The endpoint Kibana URL. Or you can use environment variable `KIBANA_URL`.
- **username**: (optional) The username to connect on it. Or you can use environment variable `KIBANA_USERNAME`.
- **password**: (optional) The password to connect on it. Or you can use environment variable `KIBANA_PASSWORD`.
- **insecure**: (optional) To disable the certificate check.
- **cacert_files**: (optional) The list of CA contend to use if you use custom PKI.
- **retry**: (optional) The number of time you should to retry connexion befaore exist with error. Default to `6`.
- **wait_before_retry**: (optional) The number of time in second we wait before each connexion retry. Default to `10`.

## Resource

- [kibana_user_space](resources/kibana_user_space.md)
- [kibana_role](resources/kibana_role.md)
- [kibana_object](resources/kibana_object.md)
- [kibana_logstash_pipeline](resources/kibana_logstash_pipeline.md)
- [kibana_copy_object](resources/kibana_copy_object.md)

## Data Source

- [kibana_host](datasources/kibana_host.md)
