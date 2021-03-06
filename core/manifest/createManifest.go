package manifest

import (
	"github.com/sandstorm/dokku-enterprise-plugin/core/utility"
	"strings"
	"fmt"
)

func CreateManifest(application string) ManifestWrapper {
	applicationConfig := utility.ExecCommand("dokku", "--quiet", "config", application)

	parsedApplicationConfig := parseConfig(applicationConfig)
	manifest := new(manifest)
	manifest.Config = make(map[string]string)

	manifestWrapper := ManifestWrapper{
		Version: 1,
		AppName: application,
		Manifest: manifest,
	}

	// Database (Mariadb)
	extractMariadb(parsedApplicationConfig, "DATABASE_URL", &manifestWrapper)
	dbIndex := 1
	for configKey, _ := range parsedApplicationConfig {
		switch {
		case strings.HasPrefix(configKey, "DOKKU_MARIADB_"):
			extractMariadb(parsedApplicationConfig, configKey, &manifestWrapper)
			dbIndex++
		}
	}

	// Remove defaults from config
	if parsedApplicationConfig["DOKKU_APP_RESTORE"] == "1" {
		delete(parsedApplicationConfig, "DOKKU_APP_RESTORE")
	}
	delete(parsedApplicationConfig, "DOKKU_PROXY_PORT_MAP")
	delete(parsedApplicationConfig, "DOKKU_DOCKERFILE_PORTS")
	delete(parsedApplicationConfig, "NO_VHOST")
	delete(parsedApplicationConfig, "DOKKU_APP_TYPE")

	// Remaining config
	for configKey, configValue := range parsedApplicationConfig {
		manifestWrapper.Manifest.Config[configKey] = configValue
	}

	// Docker options
	extractDockerOptions(&manifestWrapper, "deploy")
	extractDockerOptions(&manifestWrapper, "run")
	extractDockerOptions(&manifestWrapper, "build")

	return manifestWrapper

}

/************************
 EXTRACTORS
 */

// returns the DB name, if found!
func extractMariadb(parsedApplicationConfig map[string]string, configKey string, manifestWrapper *ManifestWrapper) string {
	dbUrl, dbUrlExists := parsedApplicationConfig[configKey]
	if dbUrlExists {
		switch {
		case strings.Contains(dbUrl, "mysql://mariadb"):
			manifestWrapper.Manifest.Mariadb = append(manifestWrapper.Manifest.Mariadb, replaceApplicationNameInString(extractDbName(dbUrl), manifestWrapper, "mariadb." + configKey, ""))

			delete(parsedApplicationConfig, configKey)
			return extractDbName(dbUrl)
		default:
			manifestWrapper.Errors = append(manifestWrapper.Errors, fmt.Sprintf("Could not parse DB URL, which was %v", dbUrl))
		}
	} else {
		manifestWrapper.DebugInfo = append(manifestWrapper.DebugInfo, "Did not find DB.")
	}
	return ""
}

func extractDockerOptions(manifestWrapper *ManifestWrapper, phase string) {
	dockerOptions := utility.ExecCommand("dokku", "docker-options", manifestWrapper.AppName, phase)

	for _, line := range strings.Split(dockerOptions, "\n") {
		line = strings.TrimSpace(line)
		switch line {
		case "--restart=on-failure:10":
			// Default value; we don't need to include this
			continue
		case "Deploy options:", "Build options:", "Run options:",
			"Deploy options: none", "Build options: none", "Run options: none":
			// first line; we can skip this.
			continue
		default:
			switch (phase) {
			case "deploy":
				manifestWrapper.Manifest.DockerOptions.Deploy = append(manifestWrapper.Manifest.DockerOptions.Deploy, replaceApplicationNameInString(line, manifestWrapper, "dockerOptions.deploy", "-v "))
			case "run":
				manifestWrapper.Manifest.DockerOptions.Run = append(manifestWrapper.Manifest.DockerOptions.Run, replaceApplicationNameInString(line, manifestWrapper, "dockerOptions.run", "-v "))
			case "build":
				manifestWrapper.Manifest.DockerOptions.Build = append(manifestWrapper.Manifest.DockerOptions.Build, replaceApplicationNameInString(line, manifestWrapper, "dockerOptions.build", "-v "))
			default:
				manifestWrapper.Errors = append(manifestWrapper.Errors, fmt.Sprintf("Unknown phase %v given; this error should never happen", phase))
			}
		}
	}
}

/************************
 HELPERS
 */
func replaceApplicationNameInString(s string, manifestWrapper *ManifestWrapper, key string, addErrorOnlyIfStringStartsWith string) string {
	r := ReplaceAppNameWithPlaceholder(s, manifestWrapper.AppName)

	if !strings.Contains(r, "[appName]") && strings.HasPrefix(s, addErrorOnlyIfStringStartsWith) {
		manifestWrapper.Errors = append(manifestWrapper.Errors, fmt.Sprintf("%v: did not find application name '%v' inside string: %v.", key, manifestWrapper.AppName, s))
	}

	return r
}
func extractDbName(dbUrl string) string {
	return dbUrl[strings.LastIndex(dbUrl, "/") + 1:]
}
func parseConfig(applicationConfig string) map[string]string {
	parsed := make(map[string]string)

	for _, line := range strings.Split(applicationConfig, "\n") {
		split := strings.SplitN(line, ":", 2)
		parsed[strings.TrimSpace(split[0])] = strings.TrimSpace(split[1])
	}
	return parsed
}