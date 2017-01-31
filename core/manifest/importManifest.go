package manifest

import (
	"github.com/sandstorm/dokku-enterprise-plugin/core/utility"
	"strings"
	"fmt"
)

func ImportManifest(application string, manifestAsString string) {
	allAppsAsString := utility.ExecCommand("dokku", "--quiet", "apps", application)
	allApps := strings.Split(allAppsAsString, "\n")

	if stringInSlice(application, allApps) {
		fmt.Printf("ERROR: Application '%v' already exists!\n", application)
		return
	}

	manifestWrapper := DeserializeManifest([]byte(manifestAsString))

	if len(manifestWrapper.Errors) > 0 {
		fmt.Printf("ERROR: The manifest had errors; which means that the manifest is NOT fully self-contained and cannot be imported: \n  %v", manifestWrapper.Errors)
		return
	}

	utility.ExecCommand("dokku", "apps:create", application)

	for _, databaseName := range manifestWrapper.Manifest.Mariadb {
		utility.ExecCommand("dokku", "mariadb:create", ReplaceAppNamePlaceholder(databaseName, application))
		utility.ExecCommand("dokku", "mariadb:link", ReplaceAppNamePlaceholder(databaseName, application), application)
	}

	for _, option := range manifestWrapper.Manifest.DockerOptions.Deploy {
		utility.ExecCommand("dokku", "docker-options:add", application, "deploy", ReplaceAppNamePlaceholder(option, application))
	}

	for _, option := range manifestWrapper.Manifest.DockerOptions.Run {
		utility.ExecCommand("dokku", "docker-options:add", application, "run", ReplaceAppNamePlaceholder(option, application))
	}

	for _, option := range manifestWrapper.Manifest.DockerOptions.Build {
		utility.ExecCommand("dokku", "docker-options:add", application, "build", ReplaceAppNamePlaceholder(option, application))
	}

	for k, v := range manifestWrapper.Manifest.Config {
		utility.ExecCommand("dokku", "config:set", application, k + "=" + ReplaceAppNamePlaceholder(v, application))
	}
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}