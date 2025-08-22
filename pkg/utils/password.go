package utils

import (
	"golang.org/x/crypto/bcrypt"
)

// HashPassword hashes a password using bcrypt
func HashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}

// CheckPasswordHash compares a password with its hash
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// ValidatePasswordStrength validates password strength
func ValidatePasswordStrength(password string) []string {
	var errors []string

	if len(password) < 6 {
		errors = append(errors, "Password must be at least 6 characters long")
	}

	if len(password) > 128 {
		errors = append(errors, "Password must be less than 128 characters")
	}

	// You can add more validation rules here
	// hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	// hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
	// hasNumber := regexp.MustCompile(`[0-9]`).MatchString(password)
	// hasSpecial := regexp.MustCompile(`[!@#$%^&*(),.?":{}|<>]`).MatchString(password)

	return errors
}
