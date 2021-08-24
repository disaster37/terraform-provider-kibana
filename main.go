package main

import (
	"os"

	"github.com/disaster37/terraform-provider-kibana/v7/kb"

	"github.com/hashicorp/terraform-plugin-sdk/plugin"
	log "github.com/sirupsen/logrus"
	"github.com/t-tomalak/logrus-easy-formatter"
)

func main() {
	log.SetFormatter(&easy.Formatter{
		LogFormat: "[%lvl%] %msg%",
	})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)

	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: kb.Provider,
	})

}
