package domain

import (
	"fmt"
	"net/mail"
	"regexp"
	"strconv"
)

var isValidName = regexp.MustCompile(`^[a-zA-Z\s]+$`).MatchString

func ValidateString(value string, minLength int, maxLength int) error {
	n := len(value)
	if n < minLength || n > maxLength {
		return fmt.Errorf("must contain from %d-%d characters", minLength, maxLength)
	}
	return nil
}

func ValidateName(value string) error {
	if err := ValidateString(value, 3, 100); err != nil {
		return err
	}
	if !isValidName(value) {
		return fmt.Errorf("must contain only lowercase letters, digits, or underscore")
	}
	return nil
}

func ValidateEmail(value string) error {
	if err := ValidateString(value, 3, 200); err != nil {
		return err
	}
	if _, err := mail.ParseAddress(value); err != nil {
		return fmt.Errorf("is not a valid email address")
	}
	return nil
}
func ValidatePassword(value string) error {
	return ValidateString(value, 6, 100)
}
func ValidateUserIdString(userId string) error {
	id, err := strconv.Atoi(userId)
	if err != nil {
		return fmt.Errorf("user ID must be a valid integer")
	}

	// Assuming a valid user ID range, adjust as needed
	if id <= 0 {
		return fmt.Errorf("user ID must be a positive integer")
	}

	// Additional string-specific checks if needed

	return nil
}
