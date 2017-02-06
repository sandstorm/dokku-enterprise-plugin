package persistentData

import (
	"github.com/sandstorm/dokku-enterprise-plugin/core/manifest"
	"strconv"
	"github.com/mholt/archiver"
	"log"
	"github.com/sandstorm/dokku-enterprise-plugin/core/utility"
	"strings"
	"os"
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
		mariadbName := manifest.ReplacePlaceholderWithAppName(mariadbNameWithPlaceholder, applicationName)

		file, err := os.Open(persistentDataDir + "/mariadb/" + strconv.Itoa(i) + ".sql")
		if err != nil {
			log.Fatalf("ERROR: could not get file handler for database file, error was: %v", err)
		}

		utility.ExecCommandWithStdIn(file,"dokku", "mariadb:import", mariadbName)

		file.Close()
	}

	// VOLUME
	allDockerOptions := manifestWrapper.GetDockerOptions()
	for _, option := range allDockerOptions {
		if (option[0:3] == "-v ") {
			volumeParts := strings.SplitN(option[3:], ":", 2)

			sourceDirectory := persistentDataDir + "/volume" + volumeParts[0]
			targetDirectory := manifest.ReplacePlaceholderWithAppName(volumeParts[0], applicationName)

			utility.CopyAndOverrideDirectory(sourceDirectory, targetDirectory)
		}
	}
}
