package persistentData

import (
	"github.com/sandstorm/dokku-enterprise-plugin/core/manifest"
	"fmt"
	"strconv"
	"github.com/mholt/archiver"
	"log"
	"github.com/sandstorm/dokku-enterprise-plugin/core/utility"
	"strings"
)

func ImportPersistentData(applicationName string, manifestWrapper manifest.ManifestWrapper, persistentDataFilePath, importTempDir string) {
	persistentDataDir := importTempDir + "/persistent-data"

	// Extracting tar.gz
	err := archiver.TarGz.Open(persistentDataFilePath, importTempDir)
	if err != nil {
		log.Fatalf("ERROR: could extract tar.gz file, error was: %v", err)
	}

	// MARIADB
	for i, mariadbNameWithPlaceholder := range manifestWrapper.Manifest.Mariadb {
		mariadbName := manifest.ReplaceAppNamePlaceholder(mariadbNameWithPlaceholder, applicationName)

		fmt.Printf("Importing database %s from %s", mariadbName, persistentDataDir + "/mariadb/" + strconv.Itoa(i) + ".sql")

		utility.ExecCommand("dokku", "mariadb:import", mariadbName, "< " + persistentDataDir + "/mariadb/" + strconv.Itoa(i) + ".sql")
	}

	// VOLUME
	allDockerOptions := manifestWrapper.GetDockerOptions()
	for _, option := range allDockerOptions {
		if (option[0:3] == "-v ") {
			volumeParts := strings.SplitN(option[3:], ":", 2)

			sourceDirectory := persistentDataDir + "/volume" + volumeParts[0]
			targetDirectory := manifest.ReplaceAppNamePlaceholder(volumeParts[0], applicationName)

			utility.CopyAndOverrideDirectory(sourceDirectory, targetDirectory)
		}
	}
}
