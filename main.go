package main

import (
	"flag"
	"math/rand"
	"mvault/credentials"
	"mvault/encryption"
	"mvault/filesys"
	"mvault/interaction"
	"mvault/network"
	"strings"
	"time"
)

/*Credentials is used to store user specific password, salt*/
type Credentials = credentials.Credentials

var path = flag.String("file", "", "path to the `file` to encrypt/decrypt")
var user = flag.String("user", "", "user name")
var local = flag.Bool("local", false, "do not interact with server")
var help = flag.Bool("help", false, "show docs")
var encr = flag.Bool("encrypt", false, "encrypt file")
var decr = flag.Bool("decrypt", false, "decrypt file")

func getSecrets(local bool) (Credentials, error) {
	var cr Credentials
	if local {
		return interaction.GetSecrets()
	}
	password, readErr := interaction.RetrieveSecret("Enter password (or leave empty): ")
	if readErr != nil {
		return cr, readErr
	}
	return network.GetCreds(*user, password)
}

func main() {
	rand.Seed(time.Now().UnixNano())
	flag.BoolVar(local, "l", *local, "alias for `-local`")
	flag.StringVar(user, "u", *user, "alias for `-user`")
	flag.BoolVar(help, "h", *help, "alias for `-help`")
	flag.BoolVar(encr, "e", *encr, "alias for `-encrypt`")
	flag.BoolVar(decr, "d", *decr, "alias for `-decrypt`")
	flag.Parse()

	if *help {
		interaction.Help("")
		return
	}

	if len(*path) < 1 || (!*encr && !*decr) || (*encr && *decr) || (len(*user) < 1 && !*local) {
		interaction.Help("Check options or filepath")
		return
	}

	var creds Credentials

	creds, getErr := getSecrets(*local)
	if getErr != nil {
		errortext := getErr.Error()
		if strings.Contains(errortext, "401") {
			interaction.Help("Wrong credentials. Retry later.")
			return
		}
		if strings.Contains(errortext, "423") {
			interaction.BanNotification()
			return
		}
		if strings.Contains(errortext, "403") {
			pwd, readErr := interaction.UpdatePassword()
			if readErr != nil {
				interaction.Help(readErr.Error())
				return
			}
			updateErr := network.Register(*user, pwd)
			if updateErr != nil {
				interaction.Help(updateErr.Error())
				return
			}

			var newErr error
			creds, newErr = network.GetCreds(*user, pwd)
			if newErr != nil {
				interaction.Help(newErr.Error())
				return
			}
		} else {
			interaction.Help(errortext)
			return
		}
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
