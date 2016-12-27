package main

import (
	"fmt"
	"github.com/sandstorm/dokku-enterprise-plugin/core/applicationLifecycleLogging"
	"github.com/sandstorm/dokku-enterprise-plugin/core/dokku"
	"os"
	"github.com/sandstorm/dokku-enterprise-plugin/core/manifest"
)

// http://dokku.viewdocs.io/dokku/development/plugin-creation/
func main() {
	command := os.Args[1]

	switch command {
	case "manifest:export":
		fmt.Println(manifest.CreateManifest(os.Args[2]))
	case "collectMetrics":
		applicationLifecycleLogging.TryToSendToServer()
		fmt.Println("Collect Metrics Done.")
	default:
		os.Exit(dokku.DokkuNotImplementedExit())
	}
}
