package configuration

import "log"

type configuration struct {
	ApiEndpointUrl string `json:"apiEndpointUrl"`
	CloudBackup cloudBackup `json:"cloudBackup"`
}

type cloudBackup struct {
	EncryptionKey string `json:"encryptionKey"`
	cloudType string `json:"cloudType"`
}

func (c cloudBackup) GetEncryptionKey() []byte {
	if len(c.EncryptionKey) < 32 {
		log.Fatalf("cloudBackup.encryptionKey must be at least 32 chars long, it was %d chars long.", len(c.EncryptionKey))
	}
	return []byte(c.EncryptionKey)
}
