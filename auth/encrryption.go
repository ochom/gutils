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

// Vault provides AES-GCM encryption and decryption functionality.
//
// Vault uses a 32-byte key for AES-256 encryption with GCM (Galois/Counter Mode)
// for authenticated encryption.
type Vault struct {
	key string
}

// NewVault creates a new Vault instance with the provided encryption key.
//
// If no key is provided, it attempts to read the key from the VAULT_KEY environment variable.
// The key must be exactly 32 bytes long for AES-256 encryption.
//
// Example:
//
//	// Using environment variable VAULT_KEY
//	vault, err := auth.NewVault()
//	if err != nil {
//		log.Fatal("Failed to create vault:", err)
//	}
//
//	// Using explicit key (must be 32 bytes)
//	key := "your-32-character-secret-key!!" // exactly 32 bytes
//	vault, err := auth.NewVault(key)
//	if err != nil {
//		log.Fatal("Invalid key length")
//	}
//
//	// Encrypt sensitive data
//	encrypted, _ := vault.Encrypt("secret password")
//	fmt.Println("Encrypted:", encrypted)
//
//	// Decrypt data
//	decrypted, _ := vault.Decrypt(encrypted)
//	fmt.Println("Decrypted:", decrypted)
func NewVault(keys ...string) (*Vault, error) {
	key := env.Get[string]("VAULT_KEY")

	if len(keys) > 0 && keys[0] != "" {
		key = keys[0]
	}

	if len(key) != 32 {
		return nil, fmt.Errorf("vault key must be 32 bytes long")
	}

	return &Vault{key: key}, nil
}

// Encrypt encrypts plaintext using AES-GCM and returns a base64-encoded string.
//
// The function generates a random nonce for each encryption, ensuring that
// encrypting the same plaintext twice produces different ciphertexts.
//
// Example:
//
//	vault, _ := auth.NewVault("your-32-character-secret-key!!")
//
//	// Encrypt a password
//	encrypted, err := vault.Encrypt("my-secret-password")
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Println("Store this:", encrypted)
//
//	// Encrypt JSON data
//	jsonData := `{"api_key": "secret123"}`
//	encryptedJSON, _ := vault.Encrypt(jsonData)
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

// Decrypt decrypts a base64-encoded ciphertext that was produced by Encrypt.
//
// Returns the original plaintext or an error if decryption fails (e.g., wrong key,
// corrupted data, or tampered ciphertext).
//
// Example:
//
//	vault, _ := auth.NewVault("your-32-character-secret-key!!")
//
//	// Decrypt previously encrypted data
//	encrypted := "base64-encoded-ciphertext..."
//	plaintext, err := vault.Decrypt(encrypted)
//	if err != nil {
//		log.Fatal("Decryption failed:", err)
//	}
//	fmt.Println("Original:", plaintext)
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
