package helpers

import (
	"fmt"
	"regexp"
)

// ParseEmail checks if the email is valid and returns the username and domain
func ParseEmail(email string) (string, string, error) {
	regex := `^([a-zA-Z0-9._%+-]+)@([a-zA-Z0-9.-]+\.[a-zA-Z]{2,})$`
	passer := regexp.MustCompile(regex)
	if !passer.MatchString(email) {
		return "", "", fmt.Errorf("invalid email")
	}

	parts := passer.FindStringSubmatch(email)
	return parts[1], parts[2], nil

}
