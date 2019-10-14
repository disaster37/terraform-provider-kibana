package main

import (
	"github.com/disaster37/terraform-provider-kibana/kb"

	"github.com/hashicorp/terraform-plugin-sdk/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: kb.Provider,
	})
}
