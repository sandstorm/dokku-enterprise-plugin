package main

import (
	"fmt"
	"github.com/sandstorm/dokku-enterprise-plugin/core/applicationLifecycleLogging"
	"github.com/sandstorm/dokku-enterprise-plugin/core/dokku"
	"os"
	"github.com/sandstorm/dokku-enterprise-plugin/core/manifest"
	"io/ioutil"
	"github.com/sandstorm/dokku-enterprise-plugin/core/cloud"
)

// http://dokku.viewdocs.io/dokku/development/plugin-creation/
func main() {
	command := os.Args[1]

	switch command {
	case "manifest:export":
		manifestWrapper := manifest.CreateManifest(os.Args[2])
		fmt.Println(string(manifest.SerializeManifest(manifestWrapper)))

	case "manifest:import":
		bytes, _ := ioutil.ReadAll(os.Stdin)
		manifest.ImportManifest(os.Args[2], string(bytes))

	case "cloud:backup":
		application := os.Args[2]
		cloud.Backup(application)

	case "cloud:createAppFromCloud":
		application := os.Args[2]
		applicationTemplate := os.Args[3]
		cloud.CreateAppFromCloud(application, applicationTemplate)

	case "collectMetrics":
		applicationLifecycleLogging.TryToSendToServer()
		fmt.Println("Collect Metrics Done.")

	default:
		os.Exit(dokku.DokkuNotImplementedExit())
	}
}
