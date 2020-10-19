package network

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

// Credentials is used to get user specific password, salt from server
type Credentials struct {
	Password string `json:"password"`
	Salt     string `json:"salt"`
}

func GetCreds(login string, password string) (string, string, error) {
	resp, err := http.Get("https://run.mocky.io/v3/9b7bfd50-241d-4e3c-8633-971fc9c5e808")
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", "", err
	}

	var creds Credentials
	err = json.Unmarshal(body, &creds)
	if err != nil {
		return "", "", err
	}

	log.Println(creds.Password)
	log.Println(creds.Salt)
	return creds.Password, creds.Salt, nil
}
