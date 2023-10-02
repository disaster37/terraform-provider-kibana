# terraform-provider-kibana

[![CircleCI](https://circleci.com/gh/disaster37/terraform-provider-kibana/tree/7.x.svg?style=svg)](https://circleci.com/gh/disaster37/terraform-provider-kibana/tree/7.x)
[![Go Report Card](https://goreportcard.com/badge/github.com/disaster37/terraform-provider-kibana)](https://goreportcard.com/report/github.com/disaster37/terraform-provider-kibana)
[![GoDoc](https://godoc.org/github.com/disaster37/terraform-provider-kibana?status.svg)](http://godoc.org/github.com/disaster37/terraform-provider-kibana)
[![codecov](https://codecov.io/gh/disaster37/terraform-provider-kibana/branch/7.x/graph/badge.svg)](https://codecov.io/gh/disaster37/terraform-provider-kibana/branch/7.x)

This is a terraform provider that lets you provision kibana resources, compatible with v7, v8 of kibana.
For Kibana 7, you need to use branch and release 7.x
For Kibana 8, you need to use branch and release 8.x

## Installation

[Download a binary](https://github.com/disaster37/terraform-provider-kibana/releases), and put it in a good spot on your system. Then update your `~/.terraformrc` to refer to the binary:

```hcl
providers {
  kibana = "/path/to/terraform-provider-kibana"
}
```

See [the docs for more information](https://www.terraform.io/docs/plugins/basics.html).


## Documentation

[Read provider documentation](docs/index.md)


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
2. Go to the right branch (7.x for Kibana 7) (`git checkout 8.x`)
3. Create your feature branch (`git checkout -b my-new-feature`)
4. Add feature, add acceptance test and test your code (`KIBANA_URL=http://127.0.0.1:5601 KIBANA_USERNAME=elastic KIBANA_PASSWORD=changeme make testacc`)
5. Commit your changes (`git commit -am 'Add some feature'`)
6. Push to the branch (`git push origin my-new-feature`)
7. Create a new Pull Request