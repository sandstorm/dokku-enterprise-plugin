package dokku

import (
	"github.com/sandstorm/dokku-enterprise-plugin/core/utility"
	"strings"
)

// Get all app container IDs as string list, or empty string if app does not exist.
func GetAppContainerIds(app string) []string {
	result := utility.ExecCommand("/bin/bash", "-c", "source $PLUGIN_CORE_AVAILABLE_PATH/common/functions; get_app_container_ids "+app)
	return strings.Split(result, " ")
}

// Get the primary app container ID (normally of the web container), or empty string if it does not exist.
func GetAppContainerId(app string) string {
	ids := GetAppContainerIds(app)
	if len(ids) > 0 {
		return ids[0]
	} else {
		return ""
	}
}

func HasAppWithName(applicationName string) bool {
	allAppsAsString := utility.ExecCommand("dokku", "--quiet", "apps")
	allApps := strings.Split(allAppsAsString, "\n")

	return utility.ContainsString(applicationName, allApps)
}

func HasMariaDBWithName(databaseName string) bool {
	allMariaDBsAsString := utility.ExecCommand("dokku", "--quiet", "mariadb:list", )
	allMariaDBs := strings.Split(allMariaDBsAsString, "\n")[1:] // ignoring the first entry which is the header

	return utility.ContainsStringStartingWith(databaseName, allMariaDBs)
}