package helpers

import (
	"crypto/sha256"
	"fmt"
	"strings"
	"unicode"

	"github.com/ochom/gutils/logs"
)

// ParseMobile  parses phone number to 254 format
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

// HashPhone hashes phone number to sha256 hash
func HashPhone(phone string) string {
	h := sha256.New()
	_, err := h.Write([]byte(phone))
	if err != nil {
		logs.Warn("could not hash the phone: %v", err)
	}

	return fmt.Sprintf("%x", h.Sum(nil))
}
