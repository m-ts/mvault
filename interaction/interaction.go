package interaction

import (
	"fmt"
	"mvault/credentials"
	"mvault/helpnote"
	"syscall"

	"golang.org/x/crypto/ssh/terminal"
)

/*Credentials is used to store user specific password, salt*/
type Credentials = credentials.Credentials

/*ACTIVE INTERACTION*/

/*RetrieveSecret asks user to write confidential info*/
func RetrieveSecret(intro string) (string, error) {
	fmt.Print(intro)
	bytePassword, readErr := terminal.ReadPassword(int(syscall.Stdin))
	if readErr != nil {
		return "", readErr
	}
	password := string(bytePassword)
	fmt.Println()
	return password, nil
}

/*GetSecrets asks user to input password and pin*/
func GetSecrets() (Credentials, error) {
	var creds Credentials

	pwd, err := RetrieveSecret("Enter password: ")
	pin, err := RetrieveSecret("Enter pin: ")
	if err != nil {
		return creds, err
	}

	creds.Password = pwd
	creds.Salt = pin

	return creds, nil
}

/*UpdatePassword updates user password*/
func UpdatePassword() (string, error) {
	fmt.Println("Your password has been reset or is not set yet.")
	pwd, err := RetrieveSecret("Enter new password (min 6): ")
	if err != nil {
		return "", err
	}
	return pwd, nil
}

/*PASSIVE INTERACTION*/

/*Help shows user guide*/
func Help(intro string) {
	if len(intro) > 0 {
		fmt.Println(intro)
		fmt.Println("For usage information type `mvault -help`")
		return
	}
	fmt.Print(helpnote.GetHelp())
}

/*BanNotification notifies user of ban event*/
func BanNotification() {
	fmt.Println("We are deeply sorry to inform you, your account has been banned and you are no logner available to use Mvault (c) for your purposes.\nIf you find these events unfair or uncalled for, please, ensure you did not violated over terms of use agreement and enduser trust code. Consider visiting psychoanalyst to find out the root of the problem.\nBest wishes,\nMvault team")
}
