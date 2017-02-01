package utility

import (
	"code.google.com/p/go.crypto/openpgp"
	"code.google.com/p/go.crypto/openpgp/armor"
	"log"
	"io/ioutil"
	"github.com/sandstorm/dokku-enterprise-plugin/core/configuration"
	"io"
)


func Encrypt(textToEncrypt io.Reader, writerForOutput io.Writer) {
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

func Decrypt(textToDecrypt io.Reader, writerForOutput io.Writer) {
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