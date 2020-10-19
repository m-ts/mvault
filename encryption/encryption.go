package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
	"log"
	"math"

	"golang.org/x/crypto/scrypt"
)

func Encrypt(data []byte, pwd string, pin string) ([]byte, []byte) {
	plaintext := data
	passwd := []byte(pwd)
	salt := []byte(pin)
	N := int(math.Pow(2, 20))
	r := 8
	p := 1
	key, err := scrypt.Key(passwd, salt, N, r, p, 32)
	// key, err := scrypt.Key(passwd, salt, 32768, 8, 1, 32)
	if err != nil {
		log.Fatal(err)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		log.Fatal(err)
	}

	nonce := make([]byte, 12)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		log.Fatal(err)
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		log.Fatal(err)
	}

	ciphertext := aesgcm.Seal(nil, nonce, plaintext, nil)
	return ciphertext, nonce
}

func Decrypt(ciphertext []byte, nonce []byte, pwd string, pin string) (string, error) {
	passwd := []byte(pwd)
	salt := []byte(pin)
	N := int(math.Pow(2, 20))
	r := 8
	p := 1
	key, err := scrypt.Key(passwd, salt, N, r, p, 32)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	plaintext, err := aesgcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}
	return string(plaintext), nil
}
