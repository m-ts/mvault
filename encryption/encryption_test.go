package encryption

import (
	"math"
	"math/rand"
	"testing"
	"time"
)

var maxdatasize = int(math.Pow(10, 5))
var mindatasize = int(math.Pow(10, 2))
var maxpwdsize = 100
var minpwdsize = 5

func randBetween(min int, max int) int {
	return rand.Intn(max-min) + min
}

/*Generates data with given size*/
func generateData(size int) []byte {
	rand.Seed(time.Now().UnixNano())
	data := make([]byte, size)
	rand.Read(data)
	return data
}

/*TestAESGCM tests createGCM function*/
func TestAESGCM(t *testing.T) {
	data := generateData(randBetween(mindatasize, maxdatasize))
	pwd := string(generateData(randBetween(minpwdsize, maxpwdsize)))
	salt := string(generateData(randBetween(minpwdsize, maxpwdsize)))
	creds := Credentials{Password: pwd, Salt: salt}

	ciphertext, err := Encrypt(data, creds)
	if err != nil || len(ciphertext) < len(data)/2 {
		t.Errorf("Couldn't encrypt: %s", err.Error())
	}

	plaintext, err := Decrypt(ciphertext, creds)
	if err != nil {
		t.Errorf("Couldn't decrypt: %s", err.Error())
	}

	if string(plaintext) != string(data) {
		t.Errorf("Decrypt(Encrypt(data, ..), ..) != data")
	}
}
