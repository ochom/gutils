package helpers

import (
	"crypto/sha256"
	"fmt"
	"strings"
	"unicode"

	"github.com/ochom/gutils/logs"
)

// ParseMobile parses a Kenyan phone number and converts it to the international format (254...).
// Returns the formatted number and a boolean indicating success.
//
// Supported input formats:
//   - 254XXXXXXXXX (international)
//   - 07XXXXXXXX or 01XXXXXXXX (local with leading 0)
//   - 7XXXXXXXX or 1XXXXXXXX (local without leading 0)
//
// Only Safaricom (7XX) and Airtel/Telkom (1XX) prefixes are supported.
//
// Example:
//
//	// Various input formats
//	phone, ok := helpers.ParseMobile("0712345678")
//	// phone = "254712345678", ok = true
//
//	phone, ok := helpers.ParseMobile("+254 712 345 678")
//	// phone = "254712345678", ok = true (spaces and + are ignored)
//
//	phone, ok := helpers.ParseMobile("712345678")
//	// phone = "254712345678", ok = true
//
//	// Invalid numbers
//	phone, ok := helpers.ParseMobile("123456")
//	// phone = "", ok = false
func ParseMobile(mobile string) (string, bool) {
	var digits []rune
	for _, r := range mobile {
		if unicode.IsDigit(r) {
			digits = append(digits, r)
		}
	}
	cleaned := string(digits)

	switch {
	case strings.HasPrefix(cleaned, "254"):
		if len(cleaned) == 12 && (cleaned[3] == '7' || cleaned[3] == '1') {
			return cleaned, true
		}
	case strings.HasPrefix(cleaned, "07") || strings.HasPrefix(cleaned, "01"):
		if len(cleaned) == 10 {
			return "254" + cleaned[1:], true
		}
	case len(cleaned) == 9 && (cleaned[0] == '7' || cleaned[0] == '1'):
		return "254" + cleaned, true
	}

	return "", false
}

// HashPhone creates a SHA-256 hash of a phone number.
// Useful for creating anonymous identifiers while maintaining the ability to lookup.
//
// Example:
//
//	// Create a hash for storage
//	hash := helpers.HashPhone("254712345678")
//	// hash = "a1b2c3..." (64-character hex string)
//
//	// Use for anonymous tracking
//	analytics.Track(helpers.HashPhone(user.Phone), event)
func HashPhone(phone string) string {
	h := sha256.New()
	_, err := h.Write([]byte(phone))
	if err != nil {
		logs.Warn("could not hash the phone: %v", err)
	}

	return fmt.Sprintf("%x", h.Sum(nil))
}
