package helpers

import "strings"

// ParseMobile ...
func ParseMobile(mobile string) string {
	// replace all non-digit characters
	mobile = strings.Map(func(r rune) rune {
		if r >= '0' && r <= '9' {
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
