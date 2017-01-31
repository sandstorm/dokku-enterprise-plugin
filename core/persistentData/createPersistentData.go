package persistentData

import (
	"github.com/sandstorm/dokku-enterprise-plugin/core/manifest"
	"strings"
	"os"
	"github.com/sandstorm/dokku-enterprise-plugin/core/utility"
	"strconv"
	"github.com/mholt/archiver"
	"log"
	"path/filepath"
	"fmt"
)

func CreatePersistentData(manifestWrapper manifest.ManifestWrapper, exportTempDir, persistentDataFilePath string) {
	persistentDataDir := exportTempDir + "/persistent-data"
	os.MkdirAll(persistentDataDir, 0777);

	// MARIADB
	os.MkdirAll(persistentDataDir + "/mariadb", 0777);
	for i, mariadbNameWithPlaceholder := range manifestWrapper.Manifest.Mariadb {
		mariadbName := replacePlaceholder(mariadbNameWithPlaceholder, manifestWrapper.AppName)
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

			targetDirectory := persistentDataDir + "/volume" + volumeParts[0]
			sourceDirectory := replacePlaceholder(volumeParts[0], manifestWrapper.AppName)

			copyDirectory(sourceDirectory, targetDirectory)
		}
	}

	// Tar.gz!
	err := archiver.TarGz.Make(persistentDataFilePath, []string{persistentDataDir})
	if err != nil {
		log.Fatalf("ERROR: could create tar.gz file, error was: %v", err)
	}
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

func ImportPersistentData(applicationName string, manifestWrapper manifest.ManifestWrapper, persistentDataFilePath, importTempDir string) {
	persistentDataDir := importTempDir + "/persistent-data"

	// Extracting tar.gz
	err := archiver.TarGz.Open(persistentDataFilePath, importTempDir)
	if err != nil {
		log.Fatalf("ERROR: could extract tar.gz file, error was: %v", err)
	}

	// MARIADB
	for i, mariadbNameWithPlaceholder := range manifestWrapper.Manifest.Mariadb {
		mariadbName := replacePlaceholder(mariadbNameWithPlaceholder, applicationName)

		fmt.Printf("Importing database %s from %s", mariadbName, persistentDataDir + "/mariadb/" + strconv.Itoa(i) + ".sql")

		utility.ExecCommand("dokku", "mariadb:import", mariadbName, "< " + persistentDataDir + "/mariadb/" + strconv.Itoa(i) + ".sql")
	}

	// VOLUME
	allDockerOptions := make([]string, 0, 20)
	allDockerOptions = append(allDockerOptions, manifestWrapper.Manifest.DockerOptions.Build...)
	allDockerOptions = append(allDockerOptions, manifestWrapper.Manifest.DockerOptions.Deploy...)
	allDockerOptions = append(allDockerOptions, manifestWrapper.Manifest.DockerOptions.Run...)
	removeDuplicates(&allDockerOptions)
	for _, option := range allDockerOptions {
		if (option[0:3] == "-v ") {
			volumeParts := strings.SplitN(option[3:], ":", 2)

			sourceDirectory := persistentDataDir + "/volume" + volumeParts[0]
			targetDirectory := replacePlaceholder(volumeParts[0], applicationName)

			copyDirectory(sourceDirectory, targetDirectory)
		}
	}

}

func replacePlaceholder(content, applicationName string) string {
	return strings.Replace(content, "[appName]", applicationName, -1)
}

func copyDirectory(source, target string) {
	source = filepath.Clean(source)
	target = filepath.Clean(target)

	os.MkdirAll(filepath.Dir(target), 0777)

	// we have to make sure that target does not exist
	// -> otherwise the complete source folder will be copied into target (instead of only its contents)
	err := os.RemoveAll(target)
	if err != nil {
		fmt.Errorf("Could not remove persistent-data-folder %s. Error was: %v", target, err)
	}

	utility.ExecCommandAndFailWithFatalErrorOnError("cp", "-R", source, target)
}