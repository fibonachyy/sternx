package utils

import "strings"

func MaskEmail(email string) string {
	// Determine the number of characters to redact (e.g., 3 characters)
	redactLength := 5

	// Ensure redactLength is within the bounds of the email length
	if redactLength > len(email) {
		redactLength = len(email)
	}

	// Implement your logic to mask or redact sensitive information
	// In this example, we redact a small portion of the email
	masked := strings.Repeat("*", redactLength)
	return masked + email[redactLength:]
}
