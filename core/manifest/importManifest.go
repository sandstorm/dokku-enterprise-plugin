package manifest

import (
	"github.com/sandstorm/dokku-enterprise-plugin/core/utility"
	"fmt"
	"github.com/sandstorm/dokku-enterprise-plugin/core/dokku"
	"strings"
)

func ImportManifest(application, manifestAsString string) {
	manifestWrapper := DeserializeManifest([]byte(manifestAsString))

	utility.ExecCommand("dokku", "apps:create", application)

	for _, databaseNameWithPlaceholder := range manifestWrapper.Manifest.Mariadb {
		databaseName := ReplacePlaceholderWithAppName(databaseNameWithPlaceholder, application)
		utility.ExecCommand("dokku", "mariadb:create", ReplacePlaceholderWithAppName(databaseName, application))
		utility.ExecCommand("dokku", "mariadb:link", ReplacePlaceholderWithAppName(databaseName, application), application)
	}

	for _, option := range manifestWrapper.Manifest.DockerOptions.Deploy {
		utility.ExecCommand("dokku", "docker-options:add", application, "deploy", ReplacePlaceholderWithAppName(option, application))
	}

	for _, option := range manifestWrapper.Manifest.DockerOptions.Run {
		utility.ExecCommand("dokku", "docker-options:add", application, "run", ReplacePlaceholderWithAppName(option, application))
	}

	for _, option := range manifestWrapper.Manifest.DockerOptions.Build {
		utility.ExecCommand("dokku", "docker-options:add", application, "build", ReplacePlaceholderWithAppName(option, application))
	}

	for k, v := range manifestWrapper.Manifest.Config {
		utility.ExecCommand("dokku", "config:set", application, k + "=" + ReplacePlaceholderWithAppName(v, application))
	}
}

func ValidateImportManifest(application, manifestAsString string) (warnings []string) {
	if dokku.HasAppWithName(application) {
		warnings = append(warnings, fmt.Sprintf("Application '%v' exists already!", application))
	}

	manifestWrapper := DeserializeManifest([]byte(manifestAsString))
	warnings = append(warnings, manifestWrapper.Errors...)

	for _, databaseNameWithPlaceholder := range manifestWrapper.Manifest.Mariadb {
		databaseName := ReplacePlaceholderWithAppName(databaseNameWithPlaceholder, application)
		if dokku.HasMariaDBWithName(databaseName) {
			warnings = append(warnings, fmt.Sprintf("Database '%v' exists already!", databaseName))
		}
	}

	allDockerOptions := manifestWrapper.GetDockerOptions()
	for _, option := range allDockerOptions {
		if (option[0:3] == "-v ") {
			volumeParts := strings.SplitN(option[3:], ":", 2)

			targetDirectory := ReplacePlaceholderWithAppName(volumeParts[0], application)

			if utility.FileExists(targetDirectory) && !utility.DirectoryIsEmpty(targetDirectory) {
				warnings = append(warnings, fmt.Sprintf("Persistent volume '%s' exists already!", targetDirectory))
			}
		}
	}

	return
}