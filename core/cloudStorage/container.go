package cloudStorage

import (
	"github.com/graymeta/stow"
	"github.com/sandstorm/dokku-enterprise-plugin/core/configuration"
	"log"
	"os"
	"bytes"
	"path/filepath"
)

var container stow.Container

func getContainer() stow.Container {
	if container == nil {
		log.Print("DEBUG: Connecting to cloud storage...")
		location, err := configuration.Get().CloudBackup.ConnectToStorage()
		if err != nil {
			log.Fatalf("ERROR: could not connect to Cloud Storage, error was: %v", err)
		}

		container, err = location.Container(configuration.Get().CloudBackup.StorageBucket)
		if err != nil {
			log.Fatalf("ERROR: did not find storage bucket '%s': %v", configuration.Get().CloudBackup.StorageBucket, err)
		}
	}

	return container
}

func UploadFile(pathAndFilename string) {
	getContainer()

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

func DownloadFile(filename, target string) string {
	getContainer()

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