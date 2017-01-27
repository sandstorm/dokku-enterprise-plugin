package cloud

import (
	"io/ioutil"
	"os"
	"time"
	"fmt"
	"github.com/sandstorm/dokku-enterprise-plugin/core/dokku"
	"github.com/sandstorm/dokku-enterprise-plugin/core/configuration"
	"github.com/sandstorm/dokku-enterprise-plugin/core/manifest"
	"github.com/sandstorm/dokku-enterprise-plugin/core/persistentData"
	"log"
	"github.com/sandstorm/dokku-enterprise-plugin/core/utility"
	"github.com/graymeta/stow"
	"path/filepath"
	"github.com/mholt/archiver"
)

func Backup(application string) {
	// BASICS
	exportTempDir, err := ioutil.TempDir(os.TempDir(), "storage-export")
	if err != nil {
		log.Fatalf("ERROR while creating temp dir: %v", err)
	}
	defer os.RemoveAll(exportTempDir)

	log.Printf("INFO: Starting export of application %s", application)
	log.Printf("DEBUG: Temp Dir: %s", exportTempDir)

	t := time.Now()
	fileBaseName := fmt.Sprintf("%s__%s__%s", application, t.Format("2006-01-02_15-04-05"), dokku.Hostname())
	filePathAndBaseName := exportTempDir + "/" + fileBaseName

	manifestFilePath := filePathAndBaseName + "-manifest.json"
	persistentDataFilePath := filePathAndBaseName + "-persistent-data.tar.gz"
	codeFilePath := filePathAndBaseName + "-code.tar.gz"

	// Cloud Connect
	log.Print("DEBUG: Uploading to cloud storage...")
	location, err := configuration.Get().CloudBackup.ConnectToStorage()
	if err != nil {
		log.Fatalf("ERROR: could not connect to Cloud Storage, error was: %v", err)
	}

	container, err := location.Container(configuration.Get().CloudBackup.StorageBucket)
	if err != nil {
		log.Fatalf("ERROR: did not find storage bucket '%s': %v", configuration.Get().CloudBackup.StorageBucket, err)
	}

	// MANIFEST
	manifestWrapper := manifest.CreateManifest(application)
	manifestBytes := manifest.SerializeManifest(manifestWrapper)
	err = ioutil.WriteFile(manifestFilePath, manifestBytes, 0755)
	log.Printf("INFO: Manifest created. Manifest is: \n%s", string(manifestBytes))

	encryptedPathAndFilename := encryptFile(manifestFilePath)
	uploadFile(encryptedPathAndFilename, container)

	// PERSISTENT DATA
	persistentData.CreatePersistentData(manifestWrapper, exportTempDir, persistentDataFilePath)
	encryptedPathAndFilename = encryptFile(persistentDataFilePath)
	uploadFile(encryptedPathAndFilename, container)

	// GIT
	err = archiver.TarGz.Make(codeFilePath, []string{
		"/home/dokku/" + application + "/config",
		"/home/dokku/" + application + "/branches",
		"/home/dokku/" + application + "/description",
		"/home/dokku/" + application + "/hooks",
		"/home/dokku/" + application + "/info",
		"/home/dokku/" + application + "/objects",
		"/home/dokku/" + application + "/refs",
	})
	if err != nil {
		log.Fatalf("ERROR: could not create tar.gz file, error was: %v", err)
	}
	encryptedPathAndFilename = encryptFile(codeFilePath)
	uploadFile(encryptedPathAndFilename, container)
}
func encryptFile(unencryptedPathAndFilename string) string {
	encryptedPathAndFilename := unencryptedPathAndFilename + ".gpg"
	gpgFile, err := os.Create(encryptedPathAndFilename)
	if err != nil {
		log.Fatalf("ERROR: %s could not be created, error was: %v", encryptedPathAndFilename, err)
	}

	sourceFile, err := os.Open(unencryptedPathAndFilename)
	if err != nil {
		log.Fatalf("ERROR: %s could not be opened, error was: %v", unencryptedPathAndFilename, err)
	}
	defer sourceFile.Close()

	log.Printf("DEBUG: encrypting %s", unencryptedPathAndFilename)
	utility.Encrypt(sourceFile, gpgFile)
	gpgFile.Close()
	return encryptedPathAndFilename
}
func uploadFile(pathAndFilename string, container stow.Container) {
	file, err := os.Open(pathAndFilename)
	if err != nil {
		log.Fatalf("ERROR: %s could not be read, error was: %v", pathAndFilename, err)
	}
	fileInfo, err := file.Stat()
	if err != nil {
		log.Fatalf("ERROR: file size for %s could not be read, error was: %v", pathAndFilename, err)
	}

	_, err = container.Put(filepath.Base(pathAndFilename), file, fileInfo.Size(), nil)
	if err != nil {
		log.Fatalf("ERROR: %s could not be uploaded, error was: %v", filepath.Base(pathAndFilename), err)
	}
}