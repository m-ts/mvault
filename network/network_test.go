package network

import (
	"math/rand"
	"strings"
	"testing"
	"time"
)

/*Generates data with given size*/
func generateData(size int) []byte {
	rand.Seed(time.Now().UnixNano())
	data := make([]byte, size)
	rand.Read(data)
	return data
}

func TestGetCredsExist(t *testing.T) {
	login, password := "Эльдар", "1234567"
	_, getErr := GetCreds(login, password)
	if getErr != nil {
		t.Errorf("Couldn't get creds: %s", getErr.Error())
	}
}

func TestGetCredsNone(t *testing.T) {
	login, password := string(generateData(150)), string(generateData(150))
	res, getErr := GetCreds(login, password)
	if getErr != nil && strings.Contains(getErr.Error(), "401") {
		return
	}
	if getErr != nil {
		t.Errorf("Wrong response: %s", getErr.Error())
	}
	t.Errorf("Wrong response: %s, %s", res.Password, res.Salt)
}
