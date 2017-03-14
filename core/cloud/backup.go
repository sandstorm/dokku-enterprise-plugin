package cloud

import (
	"io/ioutil"
	"os"
	"time"
	"fmt"
	"github.com/sandstorm/dokku-enterprise-plugin/core/dokku"
	"github.com/sandstorm/dokku-enterprise-plugin/core/manifest"
	"github.com/sandstorm/dokku-enterprise-plugin/core/persistentData"
	"github.com/sandstorm/dokku-enterprise-plugin/core/cloudStorage"
	"log"
	"github.com/sandstorm/dokku-enterprise-plugin/core/utility"
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

	// MANIFEST
	manifestWrapper := manifest.CreateManifest(application)
	manifestBytes := manifest.SerializeManifest(manifestWrapper)
	err = ioutil.WriteFile(manifestFilePath, manifestBytes, 0755)
	log.Printf("INFO: Manifest created. Manifest is: \n%s", string(manifestBytes))

	encryptedPathAndFilename := utility.EncryptFile(manifestFilePath)
	cloudStorage.UploadFile(encryptedPathAndFilename)

	// PERSISTENT DATA
	persistentData.CreatePersistentData(manifestWrapper, exportTempDir, persistentDataFilePath)
	encryptedPathAndFilename = utility.EncryptFile(persistentDataFilePath)
	cloudStorage.UploadFile(encryptedPathAndFilename)

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
	encryptedPathAndFilename = utility.EncryptFile(codeFilePath)
	cloudStorage.UploadFile(encryptedPathAndFilename)
}