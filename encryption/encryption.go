package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
	"math"
	"mvault/credentials"

	"golang.org/x/crypto/scrypt"
)

/*Credentials is used to store password, salt for encryption*/
type Credentials = credentials.Credentials

func createGCM(creds Credentials) (cipher.AEAD, error) {
	passwd := []byte(creds.Password)
	salt := []byte(creds.Salt)
	N := int(math.Pow(2, 20))
	r := 8
	p := 1
	key, generationErr := scrypt.Key(passwd, salt, N, r, p, 32)
	// key, err := scrypt.Key(passwd, salt, 32768, 8, 1, 32)
	if generationErr != nil {
		return nil, generationErr
	}

	block, cipherErr := aes.NewCipher(key)
	if cipherErr != nil {
		return nil, cipherErr
	}

	aesgcm, gcmErr := cipher.NewGCM(block)
	if gcmErr != nil {
		return nil, gcmErr
	}

	return aesgcm, nil
}

/*Encrypt given plaintext with credentials, returns nonce+ciphertext or error*/
func Encrypt(data []byte, creds Credentials) ([]byte, error) {
	aesgcm, gcmErr := createGCM(creds)
	if gcmErr != nil {
		return []byte(""), gcmErr
	}

	nonce := make([]byte, 12)
	_, nonceErr := io.ReadFull(rand.Reader, nonce)
	if nonceErr != nil {
		return []byte(""), nonceErr
	}

	ciphertext := aesgcm.Seal(nil, nonce, data, nil)
	ciphertext = append(nonce, ciphertext...)
	return ciphertext, nil
}

/*Decrypt given data (nonce+ciphertext) with credentials, returns plaintext or error*/
func Decrypt(data []byte, creds Credentials) ([]byte, error) {
	nonce, ciphertext := data[:12], data[12:]

	aesgcm, gcmErr := createGCM(creds)
	if gcmErr != nil {
		return []byte(""), gcmErr
	}

	plaintext, decryptErr := aesgcm.Open(nil, nonce, ciphertext, nil)
	if decryptErr != nil {
		return []byte(""), decryptErr
	}

	return plaintext, nil
}
