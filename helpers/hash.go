package helpers

import (
	"crypto/sha256"
	"fmt"
)

// GenerateSha256Hash generates a SHA-256 hash of the given data
func GenerateSha256Hash(data []byte) (string, error) {
	h := sha256.New()
	_, err := h.Write(data)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
