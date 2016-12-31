package utility

import (
	"code.google.com/p/go.crypto/openpgp"
	"bytes"
	"code.google.com/p/go.crypto/openpgp/armor"
	"log"
	"io/ioutil"
	"github.com/sandstorm/dokku-enterprise-plugin/core/configuration"
)


func Encrypt(encryptionText []byte) []byte {
	encryptionType := "PGP SIGNATURE"

	encbuf := bytes.NewBuffer(nil)
	w, err := armor.Encode(encbuf, encryptionType, nil)
	if err != nil {
		log.Fatal(err)
	}
	plaintext, err := openpgp.SymmetricallyEncrypt(w, configuration.Get().CloudBackup.GetEncryptionKey(), nil, nil)

	if err != nil {
		log.Fatal(err)
	}
	_, err = plaintext.Write(encryptionText)

	plaintext.Close()
	w.Close()

	return encbuf.Bytes()
}

func Decrypt(ciphertext []byte) []byte {
	decbuf := bytes.NewBuffer(ciphertext)
	result, err := armor.Decode(decbuf)
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
	return bytes
}