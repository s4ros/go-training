package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"golang.org/x/crypto/openpgp"
	"golang.org/x/crypto/openpgp/packet"
)

func main() {
	pubKey, err := os.Open("key.pub")
	if err != nil {
		log.Panic(err)
	}
	defer pubKey.Close()

	readKey := packet.NewReader(pubKey)
	entity, err := openpgp.ReadEntity(readKey)
	if err != nil {
		log.Panic(err)
	}
	toEncryptString := "This will be encrypted!"
	toEncrypt := []byte(toEncryptString)
	buff := new(bytes.Buffer)
	w, err := openpgp.Encrypt(buff, []*openpgp.Entity{entity}, nil, nil, nil)
	if err != nil {
		log.Panic(err)
	}
	_, err = w.Write(toEncrypt)
	if err != nil {
		log.Panic(err)
	}
	fmt.Println(w)
	w.Close()
	bytes, err := ioutil.ReadAll(buff)
	if err != nil {
		log.Panic(err)
	}

	encodedString := base64.StdEncoding.EncodeToString(bytes)
	fmt.Println(encodedString)
}
