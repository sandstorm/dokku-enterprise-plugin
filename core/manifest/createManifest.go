package manifest

import (
	"github.com/sandstorm/dokku-enterprise-plugin/core/utility"
	"strings"
	"encoding/json"
	"fmt"
)

func CreateManifest(application string) string {
	applicationConfig := utility.ExecCommand("dokku", "--quiet", "config", application)

	parsedApplicationConfig := parseConfig(applicationConfig)
	manifest := new(manifest)

	manifestWrapper := manifestWrapper{
		Version: 1,
		AppName: application,
		Manifest: manifest,
	}

	dbUrl, dbUrlExists := parsedApplicationConfig["DATABASE_URL"]
	if dbUrlExists {
		switch {
		case strings.Contains(dbUrl, "mysql://mariadb"):
			manifest.Mariadb = replaceApplicationNameInString(extractDbName(dbUrl), &manifestWrapper, "mariadb")
		default:
			manifestWrapper.Errors = append(manifestWrapper.Errors, fmt.Sprintf("Could not parse DB URL, which was %v", dbUrl))
		}
	} else {
		manifestWrapper.DebugInfo = append(manifestWrapper.DebugInfo, "Did not find DB.")
	}

	manifestAsBytes, _ := json.MarshalIndent(manifestWrapper, "", "  ")
	return string(manifestAsBytes)
}
func replaceApplicationNameInString(s string, manifestWrapper *manifestWrapper, key string) string {
	r := strings.Replace(s, manifestWrapper.AppName, "[appName]", -1)

	if !strings.Contains(r, "[appName]") {
		manifestWrapper.Errors = append(manifestWrapper.Errors, fmt.Sprintf("%v: did not find application name '%v' inside string: %v", key, manifestWrapper.AppName, s))
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