package persistentData

import (
	"github.com/sandstorm/dokku-enterprise-plugin/core/manifest"
	"strings"
	"os"
	"github.com/sandstorm/dokku-enterprise-plugin/core/utility"
	"strconv"
	"github.com/mholt/archiver"
	"log"
)

func CreatePersistentData(manifestWrapper manifest.ManifestWrapper, exportTempDir, persistentDataFilePath string) {
	persistentDataDir := exportTempDir + "/persistent-data"
	os.MkdirAll(persistentDataDir, 0777);

	// MARIADB
	os.MkdirAll(persistentDataDir + "/mariadb", 0777)
	for i, mariadbNameWithPlaceholder := range manifestWrapper.Manifest.Mariadb {
		mariadbName := manifest.ReplaceAppNamePlaceholder(mariadbNameWithPlaceholder, manifestWrapper.AppName)
		utility.ExecCommandAndDumpResultToFile(persistentDataDir + "/mariadb/" + strconv.Itoa(i) + ".sql", "dokku", "mariadb:export", mariadbName)
	}

	// VOLUME
	os.MkdirAll(persistentDataDir + "/volume", 0777)
	allDockerOptions := manifestWrapper.GetDockerOptions()
	for _, option := range allDockerOptions {
		if (option[0:3] == "-v ") {
			volumeParts := strings.SplitN(option[3:], ":", 2)

			targetDirectory := persistentDataDir + "/volume" + volumeParts[0]
			sourceDirectory := manifest.ReplaceAppNamePlaceholder(volumeParts[0], manifestWrapper.AppName)

			utility.CopyAndOverrideDirectory(sourceDirectory, targetDirectory)
		}
	}

	// Tar.gz!
	err := archiver.TarGz.Make(persistentDataFilePath, []string{persistentDataDir})
	if err != nil {
		log.Fatalf("ERROR: could create tar.gz file, error was: %v", err)
	}
}