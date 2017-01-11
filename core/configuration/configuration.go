package configuration

import (
	"log"
	"github.com/graymeta/stow"
	stowgs "github.com/graymeta/stow/google"
	"encoding/json"
	"github.com/graymeta/stow/s3"
	"github.com/graymeta/stow/local"
)

type configuration struct {
	ApiEndpointUrl string `json:"apiEndpointUrl"`
	CloudBackup    cloudBackup `json:"cloudBackup"`
}

type cloudBackup struct {
	EncryptionKey    string `json:"encryptionKey"`
	StorageBucket    string `json:"storageBucket"`

	CloudType        string `json:"type"`
	AwsAccessKey     string `json:"accessKey"`
	AwsSecretKey     string `json:"secretKey"`
	AwsRegion        string `json:"region"`

	GoogleProjectId  string `json:"googleProjectId"`
	GoogleConfig     interface{} `json:"googleConfig"`

	LocalStoragePath string `json:"LocalStoragePath"`
}

func (c cloudBackup) GetEncryptionKey() []byte {
	if len(c.EncryptionKey) < 32 {
		log.Fatalf("cloudBackup.encryptionKey must be at least 32 chars long, it was %d chars long.", len(c.EncryptionKey))
	}
	return []byte(c.EncryptionKey)
}

func (c cloudBackup) ConnectToStorage() (stow.Location, error) {

	switch c.CloudType {
	case "s3":
		return stow.Dial(s3.Kind, stow.ConfigMap{
			s3.ConfigAccessKeyID: c.AwsAccessKey,
			s3.ConfigSecretKey:   c.AwsSecretKey,
			s3.ConfigRegion:      c.AwsRegion,
		})
	case "google":
		googleConfig, err := json.Marshal(c.GoogleConfig)
		if err != nil {
			log.Fatal(err)
		}

		location, err := stow.Dial(stowgs.Kind, stow.ConfigMap{
			stowgs.ConfigJSON:      string(googleConfig),
			stowgs.ConfigProjectId: c.GoogleProjectId,
		})

		if location != nil {
			defer location.Close()
		}
		return location, err
	case "local":
		return stow.Dial(local.Kind, stow.ConfigMap{
			local.ConfigKeyPath: c.LocalStoragePath,
		})
	default:
		log.Fatalf("ERROR: Cloud Type %v not supported, one of 's3, google' (or 'local' for testing) must be given", c.CloudType)
	}
	return nil, nil
}
