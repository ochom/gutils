package helpers

import (
	"fmt"
	"regexp"
)

// ParseEmail validates an email address and extracts its components.
// Returns the username (local part), domain, and an error if the email is invalid.
//
// Example:
//
//	username, domain, err := helpers.ParseEmail("john.doe@example.com")
//	if err != nil {
//		fmt.Println("Invalid email")
//		return
//	}
//	fmt.Println("Username:", username) // "john.doe"
//	fmt.Println("Domain:", domain)     // "example.com"
//
//	// Use for validation
//	if _, _, err := helpers.ParseEmail(userInput); err != nil {
//		return errors.BadRequest("invalid email format")
//	}
func ParseEmail(email string) (string, string, error) {
	regex := `^([a-zA-Z0-9._%+-]+)@([a-zA-Z0-9.-]+\.[a-zA-Z]{2,})$`
	passer := regexp.MustCompile(regex)
	if !passer.MatchString(email) {
		return "", "", fmt.Errorf("invalid email")
	}

	parts := passer.FindStringSubmatch(email)
	return parts[1], parts[2], nil

}
