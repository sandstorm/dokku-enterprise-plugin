package main

import (
	"fmt"
	"github.com/sandstorm/dokku-enterprise-plugin/core/applicationLifecycleLogging"
	"github.com/sandstorm/dokku-enterprise-plugin/core/dokku"
	"os"
	"github.com/sandstorm/dokku-enterprise-plugin/core/manifest"
	"io/ioutil"
	"github.com/sandstorm/dokku-enterprise-plugin/core/utility"
)

// http://dokku.viewdocs.io/dokku/development/plugin-creation/
func main() {
	command := os.Args[1]

	switch command {
	case "manifest:export":
		manifest := manifest.CreateManifest(os.Args[2])
		fmt.Println(string(manifest))
	case "manifest:import":
		bytes, _ := ioutil.ReadAll(os.Stdin)
		manifest.ImportManifest(os.Args[2], string(bytes))
	case "manifest:exportToStorage":
		manifest := manifest.CreateManifest(os.Args[2])
		encryptedManifest := utility.Encrypt(manifest)
		fmt.Println(string(encryptedManifest))
	case "collectMetrics":
		applicationLifecycleLogging.TryToSendToServer()
		fmt.Println("Collect Metrics Done.")
	default:
		os.Exit(dokku.DokkuNotImplementedExit())
	}
}
