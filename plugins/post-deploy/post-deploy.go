package main

import (
	"github.com/sandstorm/dokku-enterprise-plugin/core/applicationLifecycleLogging"
	"github.com/sandstorm/dokku-enterprise-plugin/core/utility"
	"os"
)

// http://dokku.viewdocs.io/dokku/development/plugin-triggers/#post-deploy
func main() {
	app := os.Args[1]
	imageTag := os.Args[4]

	utility.Log("Logging successful deploy")
	applicationLifecycleLogging.AddEvent(app, "Deployment successful! (Image Tag: "+imageTag+")")
}
