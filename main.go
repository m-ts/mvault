package main

// TODO: tests

import (
	"flag"
	"math/rand"
	"mvault/credentials"
	"mvault/encryption"
	"mvault/filesys"
	"mvault/interaction"
	"mvault/network"
	"time"
)

/*Credentials is used to store user specific password, salt*/
type Credentials = credentials.Credentials

var path = flag.String("file", "", "path to the `file` to encrypt/decrypt")
var debug = flag.String("debug", "", "path to the `file` to save logs")
var local = flag.Bool("local", false, "do not interact with server")
var help = flag.Bool("help", false, "show docs")
var encr = flag.Bool("encrypt", false, "encrypt file")
var decr = flag.Bool("decrypt", false, "decrypt file")

func getSecrets(local bool) (Credentials, error) {
	if local {
		return interaction.GetSecrets()
	}
	return network.GetCreds("", "")
}

func main() {
	rand.Seed(time.Now().UnixNano())
	flag.BoolVar(local, "l", *local, "alias for `-local`")
	flag.BoolVar(help, "h", *help, "alias for `-help`")
	flag.BoolVar(encr, "e", *encr, "alias for `-encrypt`")
	flag.BoolVar(decr, "d", *decr, "alias for `-decrypt`")
	flag.Parse()

	if *help {
		interaction.Help("")
		return
	}

	if len(*path) < 1 || (!*encr && !*decr) || (*encr && *decr) {
		interaction.Help("Check options or filepath")
	}

	var creds Credentials

	creds, getErr := getSecrets(*local)
	if getErr != nil {
		interaction.Help(getErr.Error())
		return
	}

	data, readErr := filesys.ReadFile(*path)
	if readErr != nil {
		interaction.Help(readErr.Error())
		return
	}

	var newdata []byte

	if *encr {
		ciphertext, encryptErr := encryption.Encrypt(data, creds)
		if encryptErr != nil {
			interaction.Help(encryptErr.Error())
			return
		}

		newdata = ciphertext
	}
	if *decr {
		plaintext, decryptErr := encryption.Decrypt(data, creds)
		if decryptErr != nil {
			interaction.Help(decryptErr.Error())
			return
		}

		newdata = plaintext
	}

	replaceErr := filesys.Replace(*path, newdata)
	if replaceErr != nil {
		interaction.Help(replaceErr.Error())
		return
	}
}
