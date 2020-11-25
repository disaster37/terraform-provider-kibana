package main

import (
	"os"

	"github.com/disaster37/terraform-provider-kibana/v7/kb"

	"github.com/hashicorp/terraform-plugin-sdk/plugin"
	log "github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
)

func main() {
	formatter := new(prefixed.TextFormatter)
	formatter.FullTimestamp = true
	formatter.ForceFormatting = true
	log.SetFormatter(formatter)
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)

	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: kb.Provider,
	})

}
