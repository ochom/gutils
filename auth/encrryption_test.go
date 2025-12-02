package auth_test

import (
	"testing"

	"github.com/ochom/gutils/auth"
)

func TestVault_Encrypt(t *testing.T) {
	vault, err := auth.NewVault("12345678901234567890123456789012")
	if err != nil {
		t.Fatalf("failed to create vault: %v", err)
	}

	plaintext := "Hello, World!"
	encryptedText, err := vault.Encrypt(plaintext)
	if err != nil {
		t.Fatalf("encryption failed: %v", err)
	}

	if encryptedText == plaintext {
		t.Fatalf("encrypted text should not be the same as plaintext")
	}

	vault2, err := auth.NewVault("12345678901234567890123456789012")
	if err != nil {
		t.Fatalf("failed to create vault: %v", err)
	}

	decryptedText, err := vault2.Decrypt(encryptedText)
	if err != nil {
		t.Fatalf("decryption failed: %v", err)
	}

	if decryptedText != plaintext {
		t.Fatalf("decrypted text does not match original plaintext")
	}
}
