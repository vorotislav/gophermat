package crypt

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// HashPassword generates hash string based on provided password string.
func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("generate hash failed: %w", err)
	}

	return string(hash), nil
}

// CheckPassword decrypt hash and compares with provided password string.
func CheckPassword(password string, hash string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return fmt.Errorf("compare string with hash failed: %w", err)
	}

	return nil
}
