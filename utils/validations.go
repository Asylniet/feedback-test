package utils

import (
	"regexp"
	"unicode"
)

// Define validation rules for email and password
var (
	emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	// emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@kbtu.kz`)
)

func IsValidEmail(email string) bool {
	return emailRegex.MatchString(email)
}

func IsValidPassword(password string) bool {
	var (
		upp, low, num bool
		tot           uint8
	)
	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			upp = true
			tot++
		case unicode.IsLower(char):
			low = true
			tot++
		case unicode.IsNumber(char):
			num = true
			tot++
		default:
			return false
		}
	}
	if !upp || !low || !num || tot < 8 || tot > 64 {
		return false
	}
	return true
}
