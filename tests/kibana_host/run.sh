#!/bin/sh

cat <<EOT > ${HOME}/.terraformrc
provider_installation {
    filesystem_mirror {
        path    = "${PWD}/../../registry"
        include = ["registry.terraform.io/disaster37/kibana"]
    }
    direct {
        exclude = ["registry.terraform.io/disaster37/kibana"]
    }
}
EOT

rm -rf .terraform*
terraform init
TF_LOG_PROVIDER=DEBUG terraform apply