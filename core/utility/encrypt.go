package utility

import (
	"code.google.com/p/go.crypto/openpgp"
	"code.google.com/p/go.crypto/openpgp/armor"
	"log"
	"io/ioutil"
	"github.com/sandstorm/dokku-enterprise-plugin/core/configuration"
	"io"
	"strings"
	"os"
)

func EncryptFile(unencryptedPathAndFilename string) string {
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
	encrypt(sourceFile, gpgFile)
	gpgFile.Close()
	return encryptedPathAndFilename
}

func DecryptFile(encryptedPathAndFilename string) string {
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
	decrypt(sourceFile, unencryptedFile)
	defer unencryptedFile.Close()

	return unencryptedPathAndFilename
}


func encrypt(textToEncrypt io.Reader, writerForOutput io.Writer) {
	encryptionType := "PGP SIGNATURE"

	w, err := armor.Encode(writerForOutput, encryptionType, nil)
	if err != nil {
		log.Fatal(err)
	}
	plaintext, err := openpgp.SymmetricallyEncrypt(w, configuration.Get().CloudBackup.GetEncryptionKey(), nil, nil)

	if err != nil {
		log.Fatal(err)
	}

	if _, err := io.Copy(plaintext, textToEncrypt); err != nil {
		log.Fatal(err)
	}
	plaintext.Close()
	w.Close()
}

func decrypt(textToDecrypt io.Reader, writerForOutput io.Writer) {
	result, err := armor.Decode(textToDecrypt)
	if err != nil {
		log.Fatal(err)
	}

	md, err := openpgp.ReadMessage(result.Body, nil, func(keys []openpgp.Key, symmetric bool) ([]byte, error) {
		return configuration.Get().CloudBackup.GetEncryptionKey(), nil
	}, nil)
	if err != nil {
		log.Fatal(err)
	}

	bytes, err := ioutil.ReadAll(md.UnverifiedBody)
	if err != nil {
		log.Fatal(err)
	}

	_, err = writerForOutput.Write(bytes)
	if err != nil {
		log.Fatal(err)
	}
}