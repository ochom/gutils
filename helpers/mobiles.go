package helpers

import (
	"crypto/sha256"
	"fmt"
	"slices"
	"strings"

	"github.com/ochom/gutils/logs"
)

// ParseMobile  parses phone number to 254 format
func ParseMobile(mobile string) string {
	// replace all non-digit characters
	mobile = strings.Map(func(r rune) rune {
		if slices.Contains([]rune("0123456789"), r) {
			return r
		}

		return -1
	}, mobile)

	// remove leading zeros
	mobile = strings.TrimLeft(mobile, "0")

	// remove leading 254
	mobile = strings.TrimPrefix(mobile, "254")

	// check if remaining mobile is 9 digits
	if len(mobile) != 9 {
		return ""
	}

	// check if mobile starts with 7 or 1
	if mobile[0] != '7' && mobile[0] != '1' {
		return ""
	}

	return "254" + mobile
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
