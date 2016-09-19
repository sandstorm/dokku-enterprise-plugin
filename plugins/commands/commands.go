package main

import (
	"fmt"
	"github.com/sandstorm/dokku-enterprise-plugin/core/applicationLifecycleLogging"
	"github.com/sandstorm/dokku-enterprise-plugin/core/dokku"
	"os"
)

// http://dokku.viewdocs.io/dokku/development/plugin-creation/
func main() {
	command := os.Args[1]

	switch command {
	case "collectMetrics":
		applicationLifecycleLogging.TryToSendToServer()
		fmt.Println("Collect Metrics Done.")
	default:
		os.Exit(dokku.DokkuNotImplementedExit())
	}
}
