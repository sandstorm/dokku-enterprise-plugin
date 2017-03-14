package cloud

import (
	"io/ioutil"
	"os"
	"log"
	"github.com/mholt/archiver"
	"github.com/sandstorm/dokku-enterprise-plugin/core/utility"
	"github.com/sandstorm/dokku-enterprise-plugin/core/manifest"
	"github.com/sandstorm/dokku-enterprise-plugin/core/persistentData"
	"github.com/sandstorm/dokku-enterprise-plugin/core/cloudStorage"
)

func CreateAppFromCloud(application, applicationTemplate string) {
	// BASICS
	importTempDir, err := ioutil.TempDir(os.TempDir(), "storage-import")
	if err != nil {
		log.Fatalf("ERROR while creating temp dir: %v", err)
	}
	defer os.RemoveAll(importTempDir)

	log.Printf("INFO: Starting to create application %s from cloud", application)
	log.Printf("DEBUG: Temp Dir: %s", importTempDir)

	fileBasename := resolveFileBasename(applicationTemplate)

	// MANIFEST
	manifestEncryptedFilename := fileBasename + "-manifest.json.gpg"
	manifestEncryptedLocalFilePath := cloudStorage.DownloadFile(manifestEncryptedFilename, importTempDir)
	manifestLocalFilePath := utility.DecryptFile(manifestEncryptedLocalFilePath)

	manifestAsBytes, err := ioutil.ReadFile(manifestLocalFilePath)
	if err != nil {
		log.Fatalf("ERROR: could not read local manifest file, error was: %v", err)
	}
	manifestWrapper := manifest.DeserializeManifest(manifestAsBytes)

	log.Print("INFO: Validating manifest..")
	validationErrors := manifest.ValidateImportManifest(application, string(manifestAsBytes))
	if len(validationErrors) > 0 {
		utility.LogCouldNotExecuteCommand(validationErrors)
		return
	}
	log.Print("INFO: Manifest is valid. Starting to import..")

	manifest.ImportManifest(application, string(manifestAsBytes))
	log.Print("INFO: Manifest imported successfully.")

	// PERSISTENT DATA
	persistentDataEncryptedFilename := fileBasename + "-persistent-data.tar.gz.gpg"
	persistentDataEncryptedLocalFilePath := cloudStorage.DownloadFile(persistentDataEncryptedFilename, importTempDir)
	persistentDataLocalFilePath := utility.DecryptFile(persistentDataEncryptedLocalFilePath)

	persistentData.ImportPersistentData(application, manifestWrapper, persistentDataLocalFilePath, importTempDir)
	log.Print("INFO: Persistent data successfully imported.")

	// GIT
	codeEncryptedFilename := fileBasename + "-code.tar.gz.gpg"
	codeEncryptedLocalFilePath := cloudStorage.DownloadFile(codeEncryptedFilename, importTempDir)
	codeLocalFilePath := utility.DecryptFile(codeEncryptedLocalFilePath)

	err = archiver.TarGz.Open(codeLocalFilePath, "/home/dokku/" + application)
	if err != nil {
		log.Fatalf("ERROR: could not extract code from tar.gz file, error was: %v", err)
	}
	log.Print("INFO: Code successfully imported.")

	log.Printf("INFO: Successfully deployed app '%s'.", application)

}
func resolveFileBasename(applicationTemplate string) string {
	application, err := cloudStorage.GetApplication(applicationTemplate)

	if err != nil {
		utility.LogCouldNotExecuteCommand([]string{err.Error()})
	}

	return application.Versions[0].Identifier
}