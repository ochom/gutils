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

// generateRandomString creates a cryptographically secure random string
// using alphanumeric characters.
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

// GenerateHash creates a unique SHA-256 hash from the provided keys combined with the current timestamp.
// The hash is returned as a hexadecimal string.
//
// This is useful for generating unique identifiers, tokens, or cache keys.
//
// Example:
//
//	// Generate a unique hash
//	hash := helpers.GenerateHash()
//	// hash = "a1b2c3d4..." (64-character hex string)
//
//	// Generate hash from multiple inputs
//	hash := helpers.GenerateHash("user", "123", "action")
//	// hash = unique based on inputs + timestamp
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

// GenerateOTP creates a cryptographically secure numeric OTP (One-Time Password)
// of the specified length.
//
// Example:
//
//	// Generate a 6-digit OTP
//	otp := helpers.GenerateOTP(6)
//	// otp = "382947" (random 6 digits)
//
//	// Generate a 4-digit PIN
//	pin := helpers.GenerateOTP(4)
//	// pin = "5821"
//
//	// Use for verification
//	code := helpers.GenerateOTP(6)
//	storeVerificationCode(user.Email, code)
//	sendEmail(user.Email, "Your code is: "+code)
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

// HashPassword securely hashes a password using bcrypt with the default cost.
// Returns an empty string if hashing fails.
//
// Example:
//
//	// Hash a password before storing
//	hashedPassword := helpers.HashPassword("userSecretPassword")
//	user.PasswordHash = hashedPassword
//	db.Save(user)
//
//	// The hash is safe to store in a database
//	fmt.Println(len(hashedPassword)) // ~60 characters
func HashPassword(password string) string {
	bsp, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return ""
	}

	return string(bsp)
}

// ComparePassword securely compares a bcrypt-hashed password with a plaintext password.
// Returns true if the password matches the hash, false otherwise.
//
// Example:
//
//	// Verify login credentials
//	user, _ := db.FindUserByEmail(email)
//	if !helpers.ComparePassword(user.PasswordHash, inputPassword) {
//		return errors.Unauthorized("invalid credentials")
//	}
//
//	// Password change verification
//	if !helpers.ComparePassword(user.PasswordHash, oldPassword) {
//		return errors.BadRequest("current password is incorrect")
//	}
func ComparePassword(hashedString string, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedString), []byte(password))
	return err == nil
}
