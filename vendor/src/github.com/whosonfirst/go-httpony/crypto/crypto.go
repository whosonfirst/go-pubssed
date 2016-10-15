package crypto

// https://golang.org/pkg/crypto/cipher/
// https://gist.github.com/kkirsche/e28da6754c39d5e7ea10

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"io"
	"strings"
)

type Crypt struct {
	block cipher.Block
}

func NewCrypt(key string) (*Crypt, error) {

	block, err := aes.NewCipher([]byte(key))

	if err != nil {
		return nil, err
	}

	c := Crypt{block}
	return &c, nil
}

func (c *Crypt) Encrypt(plaintext string) (string, error) {

	nonce := make([]byte, 12)

	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	aesgcm, err := cipher.NewGCM(c.block)

	if err != nil {
		return "", err
	}

	ciphertext := aesgcm.Seal(nil, nonce, []byte(plaintext), nil)

	hex_cipher := hex.EncodeToString(ciphertext)
	hex_nonce := hex.EncodeToString(nonce)

	parts := []string{hex_cipher, hex_nonce}
	return strings.Join(parts, "#"), nil
}

func (c *Crypt) Decrypt(enc string) (string, error) {

	parts := strings.Split(enc, "#")

	if len(parts) != 2 {
	   return "", errors.New("Unable to parse encrypted value")
	}
	
	hex_cipher := parts[0]
	hex_nonce := parts[1]

	ciphertext, err := hex.DecodeString(hex_cipher)

	if err != nil {
		return "", err
	}

	nonce, err := hex.DecodeString(hex_nonce)

	if err != nil {
		return "", err
	}

	aesgcm, err := cipher.NewGCM(c.block)

	if err != nil {
		return "", err
	}

	plaintext, err := aesgcm.Open(nil, nonce, ciphertext, nil)

	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}
