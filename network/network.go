package network

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"mvault/credentials"
	"net/http"
	"strconv"
)

// https://run.mocky.io/v3/9b7bfd50-241d-4e3c-8633-971fc9c5e808
// http://84.38.181.216/api/crypt

// TODO: sending user creds to get encryption creds

/*Credentials is used to get user specific password, salt from server*/
type Credentials = credentials.Credentials

/*RequestCreds is used to send user login and password*/
type RequestCreds = credentials.RequestCredentials

/*ResponseJSON contains message from server*/
type ResponseJSON struct {
	Message string `json:"message"`
}

/*GetCreds retrieves `Credentials` from server*/
func GetCreds(login string, password string) (Credentials, error) {
	var creds Credentials
	var rCreds RequestCreds
	rCreds.Login, rCreds.Password = login, password
	credsM, marshErr := json.Marshal(rCreds)
	if marshErr != nil {
		return creds, marshErr
	}
	content := bytes.NewBuffer(credsM)

	resp, postErr := http.Post("http://84.38.181.216/api/crypt", "application/json", content)
	if postErr != nil {
		return creds, postErr
	}
	defer resp.Body.Close()

	if resp.StatusCode == 423 {
		return creds, errors.New(strconv.Itoa(resp.StatusCode) + " BAN")
	}

	if resp.StatusCode != 201 {
		return creds, errors.New(strconv.Itoa(resp.StatusCode) + " status code")
	}

	body, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		return creds, readErr
	}

	unmErr := json.Unmarshal(body, &creds)
	if unmErr != nil {
		return creds, unmErr
	}

	// log.Println(creds.Password)
	// log.Println(creds.Salt)
	return creds, nil
}

/*Register sets user password*/
func Register(login string, password string) error {
	var rCreds RequestCreds
	rCreds.Login, rCreds.Password = login, password
	credsM, marshErr := json.Marshal(rCreds)
	if marshErr != nil {
		return marshErr
	}
	content := bytes.NewBuffer(credsM)

	resp, postErr := http.Post("http://84.38.181.216/api/auth/register", "application/json", content)
	if postErr != nil {
		return postErr
	}
	defer resp.Body.Close()

	body, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		return readErr
	}
	var response ResponseJSON
	unmErr := json.Unmarshal(body, &response)
	if unmErr != nil {
		return unmErr
	}

	if resp.StatusCode == 201 {
		return nil
	}

	return errors.New(response.Message)
}
