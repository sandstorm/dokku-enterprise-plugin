package main

import (
	"fmt"
	"github.com/sandstorm/dokku-enterprise-plugin/core/applicationLifecycleLogging"
	"github.com/sandstorm/dokku-enterprise-plugin/core/dokku"
	"os"
	"github.com/sandstorm/dokku-enterprise-plugin/core/manifest"
	"io/ioutil"
	"github.com/sandstorm/dokku-enterprise-plugin/core/utility"
	"github.com/sandstorm/dokku-enterprise-plugin/core/configuration"
	"log"
	"time"
	"github.com/mholt/archiver"
)

// http://dokku.viewdocs.io/dokku/development/plugin-creation/
func main() {
	command := os.Args[1]

	switch command {
	case "manifest:export":
		manifest := manifest.CreateManifest(os.Args[2])
		fmt.Println(string(manifest))
	case "manifest:import":
		bytes, _ := ioutil.ReadAll(os.Stdin)
		manifest.ImportManifest(os.Args[2], string(bytes))
	case "manifest:exportWithDataToStorage":
		application := os.Args[2]

		exportTempDir, err := ioutil.TempDir(os.TempDir(), "storage-export")
		if err != nil {
			log.Fatalf("ERROR while creating temp dir: %v", err)
		}

		log.Printf("INFO: Starting export of application %s", application)
		log.Printf("DEBUG: Temp Dir: %s", exportTempDir)

		manifest.CreateManifestAndStoreDataIntoTemporaryFolder(application, exportTempDir)
		manifestBytes, _ := ioutil.ReadFile(exportTempDir + "/manifest.json")
		log.Printf("INFO: Manifest created. Manifest is: \n%s", string(manifestBytes))

		t := time.Now()
		fileName := fmt.Sprintf("%s__%s__%s.tar.gz", application, t.Format("2006-01-02_15-04-05"), dokku.Hostname())
		tarGzFile := os.TempDir() + "/" + fileName
		log.Printf("DEBUG: exporting tar.gz to %s", tarGzFile)

		err = archiver.TarGz.Make(tarGzFile, []string{exportTempDir})
		if err != nil {
			log.Fatalf("ERROR: could create tar.gz file, error was: %v", err)
		}

		unencryptedManifest, err := os.Open(tarGzFile)
		if err != nil {
			log.Fatalf("ERROR: %s could not be opened, error was: %v", tarGzFile, err)
		}
		defer unencryptedManifest.Close()

		gpgFile, err := os.Create(tarGzFile + ".gpg")
		if err != nil {
			log.Fatalf("ERROR: %s could not be created, error was: %v", tarGzFile + ".gpg", err)
		}

		log.Printf("DEBUG: encrypting tar.gz to %s", tarGzFile + ".gpg")
		utility.Encrypt(unencryptedManifest, gpgFile)
		defer gpgFile.Close()


		log.Print("DEBUG: Uploading to cloud storage...")
		location, err := configuration.Get().CloudBackup.ConnectToStorage()
		if err != nil {
			log.Fatalf("ERROR: could not connect to Cloud Storage, error was: %v", err)
		}

		container, err := location.Container(configuration.Get().CloudBackup.GoogleStorageBucket)
		if err != nil {
			log.Fatalf("ERROR: did not find storage bucket '%s': %v", configuration.Get().CloudBackup.GoogleStorageBucket, err)
		}

		// HINT: for the google implmentation, we can IGNORE the size value! :-) (we don't have it due to streaming!)
		gpgFileForReading, err := os.Open(tarGzFile + ".gpg")
		if err != nil {
			log.Fatalf("ERROR: %s could not be read, error was: %v", tarGzFile + ".gpg", err)
		}
		fileInfo, err := gpgFileForReading.Stat()
		if err != nil {
			log.Fatalf("ERROR: file size for %s could not be read, error was: %v", tarGzFile + ".gpg", err)
		}

		container.Put(fileName + ".gpg", gpgFileForReading, fileInfo.Size(), nil)
		log.Printf("INFO: DONE! %s", fileName + ".gpg")
	case "collectMetrics":
		applicationLifecycleLogging.TryToSendToServer()
		fmt.Println("Collect Metrics Done.")
	default:
		os.Exit(dokku.DokkuNotImplementedExit())
	}
}
