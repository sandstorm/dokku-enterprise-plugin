package manifest

import (
	"github.com/sandstorm/dokku-enterprise-plugin/core/utility"
	"strings"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
)

func CreateManifest(application string) []byte {
	return CreateManifestAndStoreDataIntoTemporaryFolder(application, "")
}
/**
 * Second argument might be NULL
 */
func CreateManifestAndStoreDataIntoTemporaryFolder(application string, temporaryFolder string) []byte {
	applicationConfig := utility.ExecCommand("dokku", "--quiet", "config", application)

	parsedApplicationConfig := parseConfig(applicationConfig)
	manifest := new(manifest)
	manifest.Config = make(map[string]string)

	manifestWrapper := manifestWrapper{
		Version: 1,
		AppName: application,
		Manifest: manifest,
	}

	// Database (Mariadb)
	dbName := extractMariadb(parsedApplicationConfig, "DATABASE_URL", &manifestWrapper)
	if len(dbName) > 0 && len(temporaryFolder) > 0 {
		utility.ExecCommandAndDumpResultToFile(temporaryFolder + "/mariadb/0.sql", "dokku", "mariadb:export", dbName)
	}
	dbIndex := 1
	for configKey, _ := range parsedApplicationConfig {
		switch {
		case strings.HasPrefix(configKey, "DOKKU_MARIADB_"):
			dbName := extractMariadb(parsedApplicationConfig, configKey, &manifestWrapper)
			if len(dbName) > 0 && len(temporaryFolder) > 0 {
				utility.ExecCommandAndDumpResultToFile(temporaryFolder + "/mariadb/" + strconv.Itoa(dbIndex) + ".sql", "dokku", "mariadb:export", dbName)
			}

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

	if len(temporaryFolder) > 0 {
		allDockerOptions := make([]string, 0, 20)
		allDockerOptions = append(allDockerOptions, manifestWrapper.Manifest.DockerOptions.Build...)
		allDockerOptions = append(allDockerOptions, manifestWrapper.Manifest.DockerOptions.Run...)
		allDockerOptions = append(allDockerOptions, manifestWrapper.Manifest.DockerOptions.Deploy...)
	}

	manifestAsBytes, err := json.MarshalIndent(manifestWrapper, "", "  ")

	if err != nil {
		log.Fatalf("There was an error serializing JSON manifest: %v", err)
	}

	if len(temporaryFolder) > 0 {
		ioutil.WriteFile(temporaryFolder + "/manifest.json", manifestAsBytes, 0644)
	}
	return manifestAsBytes
}

/************************
 EXTRACTORS
 */

// returns the DB name, if found!
func extractMariadb(parsedApplicationConfig map[string]string, configKey string, manifestWrapper *manifestWrapper) string {
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

func extractDockerOptions(manifestWrapper *manifestWrapper, phase string) {
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
func replaceApplicationNameInString(s string, manifestWrapper *manifestWrapper, key string, addErrorOnlyIfStringStartsWith string) string {
	r := strings.Replace(s, manifestWrapper.AppName, "[appName]", -1)

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

func removeDuplicates(xs *[]string) {
	found := make(map[string]bool)
	j := 0
	for i, x := range *xs {
		if !found[x] {
			found[x] = true
			(*xs)[j] = (*xs)[i]
			j++
		}
	}
	*xs = (*xs)[:j]
}