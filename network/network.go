package network

import (
	"encoding/json"
	"io/ioutil"
	"mvault/credentials"
	"net/http"
)

// Credentials is used to get user specific password, salt from server
type Credentials = credentials.Credentials

/*GetCreds retrieves `Credentials` from server*/
func GetCreds(login string, password string) (Credentials, error) {
	var creds Credentials
	resp, getErr := http.Get("https://run.mocky.io/v3/9b7bfd50-241d-4e3c-8633-971fc9c5e808")
	if getErr != nil {
		return creds, getErr
	}
	defer resp.Body.Close()

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
