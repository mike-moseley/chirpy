package auth

import (
	"fmt"

	"github.com/alexedwards/argon2id"
)

func HashPassword(password string) (string, error) {
	hashedPass, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		return "", fmt.Errorf("Error hashing password: %w", err)
	}

	return hashedPass, nil
}

func CheckPasswordHash(password, hash string) (bool, error) {
	res, err := argon2id.ComparePasswordAndHash(password, hash)
	if err != nil {
		return false, fmt.Errorf("Error checking password: %w", err)
	}
	return res, nil
}
