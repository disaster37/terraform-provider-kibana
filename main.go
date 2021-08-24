package main

import (
	"context"
	"flag"
	"log"

	"github.com/disaster37/terraform-provider-kibana/v7/kb"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
)

func main() {

	var debugMode bool

	flag.BoolVar(&debugMode, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	opts := &plugin.ServeOpts{ProviderFunc: kb.Provider}

	if debugMode {
		err := plugin.Debug(context.Background(), "registry.terraform.io/disaster37/kibana", opts)

		if err != nil {
			log.Fatal(err.Error())
		}

		return
	}

	plugin.Serve(opts)

}
