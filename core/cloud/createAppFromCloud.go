package cloud

import (
	"io/ioutil"
	"os"
	"github.com/sandstorm/dokku-enterprise-plugin/core/configuration"
	"log"
	"github.com/mholt/archiver"
	"github.com/graymeta/stow"
	"bytes"
	"github.com/sandstorm/dokku-enterprise-plugin/core/utility"
	"strings"
	"github.com/sandstorm/dokku-enterprise-plugin/core/manifest"
	"github.com/sandstorm/dokku-enterprise-plugin/core/persistentData"
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

	// Cloud Connect
	log.Print("DEBUG: Connecting with cloud storage...")
	location, err := configuration.Get().CloudBackup.ConnectToStorage()
	if err != nil {
		log.Fatalf("ERROR: could not connect to Cloud Storage, error was: %v", err)
	}

	container, err := location.Container(configuration.Get().CloudBackup.StorageBucket)
	if err != nil {
		log.Fatalf("ERROR: did not find storage bucket '%s': %v", configuration.Get().CloudBackup.StorageBucket, err)
	}

	fileBasename := resolveFileBasename(applicationTemplate, container)

	// MANIFEST
	manifestEncryptedFilename := fileBasename + "-manifest.json.gpg"
	manifestEncryptedLocalFilePath := downloadFile(manifestEncryptedFilename, importTempDir, container)
	manifestLocalFilePath := decryptFile(manifestEncryptedLocalFilePath)

	manifestAsBytes, err := ioutil.ReadFile(manifestLocalFilePath)
	if err != nil {
		log.Fatalf("ERROR: could not read local manifest file, error was: %v", err)
	}

	manifest.ImportManifest(application, string(manifestAsBytes))
	log.Print("INFO: Manifest imported successfully.")

	// PERSISTENT DATA
	persistentDataEncryptedFilename := fileBasename + "-persistent-data.tar.gz.gpg"
	persistentDataEncryptedLocalFilePath := downloadFile(persistentDataEncryptedFilename, importTempDir, container)
	persistentDataLocalFilePath := decryptFile(persistentDataEncryptedLocalFilePath)

	persistentData.ImportPersistentData(application, manifest.DeserializeManifest(manifestAsBytes), persistentDataLocalFilePath, importTempDir)
	log.Print("INFO: Persistent data imported successfully.")

	// GIT
	codeEncryptedFilename := fileBasename + "-code.tar.gz.gpg"
	codeEncryptedLocalFilePath := downloadFile(codeEncryptedFilename, importTempDir, container)
	codeLocalFilePath := decryptFile(codeEncryptedLocalFilePath)

	err = archiver.TarGz.Open(codeLocalFilePath, "/home/dokku/" + application)
	if err != nil {
		log.Fatalf("ERROR: could not extract code from tar.gz file, error was: %v", err)
	}

}
func resolveFileBasename(applicationTemplate string, container stow.Container) string {
	application := getApplication(applicationTemplate, container)

	if len(application.Versions) == 0 || application.Versions[0].Identifier == "" {
		log.Fatalf("ERROR: Could not find any valid application with name %s", applicationTemplate)
	}

	return application.Versions[0].Identifier
}
func decryptFile(encryptedPathAndFilename string) string {
	unencryptedPathAndFilename := strings.TrimSuffix(encryptedPathAndFilename, ".gpg")
	unencryptedFile, err := os.Create(unencryptedPathAndFilename)
	if err != nil {
		log.Fatalf("ERROR: %s could not be created, error was: %v", unencryptedPathAndFilename, err)
	}

	sourceFile, err := os.Open(encryptedPathAndFilename)
	if err != nil {
		log.Fatalf("ERROR: %s could not be opened, error was: %v", encryptedPathAndFilename, err)
	}
	defer sourceFile.Close()

	log.Printf("DEBUG: decrypting %s", encryptedPathAndFilename)
	utility.Decrypt(sourceFile, unencryptedFile)
	defer unencryptedFile.Close()

	return unencryptedPathAndFilename
}
func downloadFile(filename, target string, container stow.Container) string {
	item, err := container.Item(container.ID() + "/" + filename)
	if err != nil {
		log.Fatalf("ERROR: could not locate item %s in cloud storage, error was: %v", filename, err)
	}

	itemReader, err := item.Open()
	if err != nil {
		log.Fatalf("ERROR: could not read item %s from cloud storage, error was: %v", filename, err)
	}
	defer itemReader.Close()

	localFilePath := target + "/" + filename
	file, err := os.Create(localFilePath)
	if err != nil {
		log.Fatalf("ERROR: could not create file at %s, error was: %v", localFilePath, err)
	}

	contentBuffer := new(bytes.Buffer)
	_, err = contentBuffer.ReadFrom(itemReader)
	if err != nil {
		log.Fatalf("ERROR: error while reading from item %s of cloud storage, error was: %v", filename, err)
	}

	file.Write(contentBuffer.Bytes())
	file.Close()

	return localFilePath
}