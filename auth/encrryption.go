package auth

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"

	"github.com/ochom/gutils/env"
)

// Vault
type Vault struct {
	key string
}

// NewVault creates a new Vault with the given key
func NewVault(keys ...string) (*Vault, error) {
	key := env.Get("VAULT_KEY")

	if len(keys) > 0 && keys[0] != "" {
		key = keys[0]
	}

	if len(key) != 32 {
		return nil, fmt.Errorf("vault key must be 32 bytes long")
	}

	return &Vault{key: key}, nil
}

// Encrypt text using AES-GCM
func (v Vault) Encrypt(plaintext string) (string, error) {
	block, err := aes.NewCipher([]byte(v.key))
	if err != nil {
		return "", err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, aesGCM.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := aesGCM.Seal(nonce, nonce, []byte(plaintext), nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// Decrypt text using AES-GCM
func (v Vault) Decrypt(encryptedText string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(encryptedText)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher([]byte(v.key))
	if err != nil {
		return "", err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := aesGCM.NonceSize()
	if len(data) < nonceSize {
		return "", fmt.Errorf("ciphertext too short")
	}

	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}
