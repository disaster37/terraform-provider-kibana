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
