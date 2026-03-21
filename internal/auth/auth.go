package auth

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/alexedwards/argon2id"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
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

func MakeJWT(userID uuid.UUID, tokenSecret string) (string, error) {
	claims := jwt.RegisteredClaims{
		Issuer:    "chirpy-access",
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		Subject:   userID.String(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedString, err := token.SignedString([]byte(tokenSecret))
	if err != nil {
		return "", fmt.Errorf("Error signing:\n  %w", err)
	}

	return signedString, nil
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	claims := &jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(tokenSecret), nil
	})
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("Error validating:\n  %w", err)
	}
	subject, err := token.Claims.GetSubject()
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("Error getting subject value:\n  %w", err)
	}
	userID, err := uuid.Parse(subject)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("Error parsing UUID:\n  %w", err)
	}
	return userID, nil
}

func GetBearerToken(headers http.Header) (string, error) {
	auth := headers.Get("Authorization")

	if auth == "" {
		return "", fmt.Errorf("No auth token received")
	}
	clean_auth, _ := strings.CutPrefix(auth, "Bearer ")
	return clean_auth, nil
}

func MakeRefreshToken() string {
	tokenBytes := make([]byte, 32)
	rand.Read(tokenBytes)
	return hex.EncodeToString(tokenBytes)
}

func GetAPIKey(headers http.Header) (string, error) {
	headerApiKey := headers.Get("Authorization")
	if headerApiKey == "" {
		return "", fmt.Errorf("No API key provided")
	}
	headerApiKey, found := strings.CutPrefix(headerApiKey, "ApiKey ")
	if found == false {
		return "", fmt.Errorf("API key malformed")
	}

	return strings.TrimSpace(headerApiKey), nil
}
