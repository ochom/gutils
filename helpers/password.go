package helpers

import (
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"math/big"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// generateRandomString generates a random string of a given size
func generateRandomString(size int) string {
	var letterRunes = []rune("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, size)

	for i := range b {
		index := func() int {
			bigN, err := rand.Int(rand.Reader, big.NewInt(int64(len(letterRunes))))
			if err != nil {
				return 0
			}
			return int(bigN.Int64())
		}()
		b[i] = letterRunes[index]
	}

	return string(b)
}

// GenerateHash creates a hash from a string
func GenerateHash(keys ...string) string {
	now := fmt.Sprintf("%d", time.Now().UnixNano())
	keys = append(keys, now)
	key := strings.Join(keys, "")
	hash := sha256.New()
	n, err := hash.Write([]byte(key))
	if err != nil || n != len(key) {
		return generateRandomString(sha256.New().Size())
	}

	return fmt.Sprintf("%x", hash.Sum(nil))
}

// GenerateOTP generates a random  OTP of a given size
func GenerateOTP(size int) string {
	var letterRunes = []rune("0123456789")
	b := make([]rune, size)

	for i := range b {
		index := func() int {
			bigN, err := rand.Int(rand.Reader, big.NewInt(int64(len(letterRunes))))
			if err != nil {
				return 0
			}
			return int(bigN.Int64())
		}()
		b[i] = letterRunes[index]
	}

	return string(b)
}

// HashPassword hashes a password
func HashPassword(password string) []byte {
	bsp, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return []byte{}
	}

	return bsp
}

// ComparePassword compares a password with a hash
func ComparePassword(hash []byte, password string) bool {
	err := bcrypt.CompareHashAndPassword(hash, []byte(password))
	return err == nil
}
