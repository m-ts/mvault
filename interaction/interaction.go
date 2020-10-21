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

func retrievesecret(intro string) (string, error) {
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

	pwd, err := retrievesecret("Enter password: ")
	pin, err := retrievesecret("Enter pin: ")
	if err != nil {
		return creds, err
	}

	creds.Password = pwd
	creds.Salt = pin

	return creds, nil
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
