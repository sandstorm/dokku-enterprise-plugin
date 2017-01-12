package persistentData

import (
	"github.com/sandstorm/dokku-enterprise-plugin/core/manifest"
	"strings"
	"os"
	"github.com/sandstorm/dokku-enterprise-plugin/core/utility"
	"strconv"
	"github.com/mholt/archiver"
	"log"
	"fmt"
)

func CreatePersistentData(manifestWrapper manifest.ManifestWrapper, exportTempDir, persistentDataFilePath string) {
	persistentDataDir := exportTempDir + "/persistent-data"
	os.MkdirAll(persistentDataDir, 0777);

	// MARIADB
	os.MkdirAll(persistentDataDir + "/mariadb", 0777);
	for i, mariadbNameWithPlaceholder := range manifestWrapper.Manifest.Mariadb {
		mariadbName := replacePlaceholder(mariadbNameWithPlaceholder, manifestWrapper)
		utility.ExecCommandAndDumpResultToFile(persistentDataDir + "/mariadb/" + strconv.Itoa(i) + ".sql", "dokku", "mariadb:export", mariadbName)
	}

	// VOLUME
	os.MkdirAll(persistentDataDir + "/volume", 0777);
	allDockerOptions := make([]string, 0, 20)
	allDockerOptions = append(allDockerOptions, manifestWrapper.Manifest.DockerOptions.Build...)
	allDockerOptions = append(allDockerOptions, manifestWrapper.Manifest.DockerOptions.Deploy...)
	allDockerOptions = append(allDockerOptions, manifestWrapper.Manifest.DockerOptions.Run...)
	removeDuplicates(&allDockerOptions)
	for _, option := range allDockerOptions {
		if (option[0:3] == "-v ") {
			volumeParts := strings.SplitN(option[3:], ":", 2)

			targetDirectory := persistentDataDir + "/volume/" + volumeParts[0]
			sourceDirectory := replacePlaceholder(volumeParts[0], manifestWrapper)
			os.MkdirAll(targetDirectory, 0777)
			utility.ExecCommand("cp", "-R", sourceDirectory, targetDirectory)
		}
	}

	// Tar.gz!
	err := archiver.TarGz.Make(persistentDataFilePath, []string{persistentDataDir})
	if err != nil {
		log.Fatalf("ERROR: could create tar.gz file, error was: %v", err)
	}
}

func replacePlaceholder(s string, wrapper manifest.ManifestWrapper) string {
	return strings.Replace(s, "[appName]", wrapper.AppName, -1)
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