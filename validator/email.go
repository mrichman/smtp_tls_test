package validator

import (
	"errors"
	"regexp"
	"strings"
)

var (
	emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
)

// ValidateEmail checks if the provided email address is valid
func ValidateEmail(email string) error {
	if email == "" {
		return errors.New("email cannot be empty")
	}

	email = strings.TrimSpace(email)
	if !emailRegex.MatchString(email) {
		return errors.New("invalid email format")
	}

	return nil
}

// ValidateEmails checks if all provided email addresses are valid
func ValidateEmails(emails []string) error {
	if len(emails) == 0 {
		return errors.New("at least one email address is required")
	}

	for _, email := range emails {
		if err := ValidateEmail(email); err != nil {
			return err
		}
	}

	return nil
}
