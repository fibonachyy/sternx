package service

import (
	"errors"
	"regexp"

	"github.com/fibonachyy/sternx/internal/api/user"
)

func validateUserRequest(req *user.UserRequest) error {
	if req.Name == "" {
		return errors.New("name cannot be empty")
	}
	if req.Email == "" {
		return errors.New("email cannot be empty")
	}
	if !isValidEmail(req.Email) {
		return errors.New("invalid email format")
	}
	return nil
}

func isValidEmail(email string) bool {
	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	match, _ := regexp.MatchString(emailRegex, email)
	return match
}
