package main

import (
	"flag"
	"os"

	"github.com/disaster37/terraform-provider-kibana/v8/kb"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	log "github.com/sirupsen/logrus"
	easy "github.com/t-tomalak/logrus-easy-formatter"
)

func init() {

	log.SetOutput(os.Stderr)
	log.SetLevel(log.DebugLevel)
	log.SetFormatter(&easy.Formatter{
		LogFormat: "[%lvl%] %msg%\n",
	})

}

func main() {

	var debugMode bool

	flag.BoolVar(&debugMode, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	opts := &plugin.ServeOpts{
		ProviderFunc: kb.Provider,
		Debug:        debugMode,
	}

	plugin.Serve(opts)

}
