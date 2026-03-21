package auth

import (
	"net/http/httptest"
	"testing"
	"time"

	"github.com/alexedwards/argon2id"
	"github.com/google/uuid"
)

func TestHashPassword(t *testing.T) {
	pass := "ver1table-smorg4sb0rd"
	hash, err := argon2id.CreateHash(pass, argon2id.DefaultParams)
	if err != nil {
		t.Error(err)
	}
	res, err := argon2id.ComparePasswordAndHash(pass, hash)
	if !res || err != nil {
		t.Errorf(`Password did not match hash`)
	}
}

func TestHashPasswordFail(t *testing.T) {
	pass := "ver1table-smorg4sb0rd"
	hash, err := argon2id.CreateHash(pass, argon2id.DefaultParams)
	if err != nil {
		t.Error(err)
	}
	hash = hash + "WRONG!"
	res, err := argon2id.ComparePasswordAndHash(pass, hash)
	if res || err == nil {
		t.Errorf(`Password did not match hash`)
	}
}

func TestHashPasswordBlank(t *testing.T) {
	pass := ""
	hash, err := argon2id.CreateHash(pass, argon2id.DefaultParams)
	if err != nil {
		t.Error(err)
	}
	res, err := argon2id.ComparePasswordAndHash(pass, hash)
	if !res || err != nil {
		t.Errorf(`Password did not match hash`)
	}
}

func TestCheckPasswordHash(t *testing.T) {
	pass := "ver1table-smorg4sb0rd"
	hash, err := argon2id.CreateHash(pass, argon2id.DefaultParams)
	if err != nil {
		t.Error(err)
	}
	res, err := CheckPasswordHash(pass, hash)
	if !res || err != nil {
		t.Errorf(`Password did not match hash`)
	}
}

func TestJWT(t *testing.T) {
	userID := uuid.New()
	tokenSecret := "0rgasbord-org4asbord"
	jwt, err := MakeJWT(userID, tokenSecret)
	if err != nil {
		t.Error(err)
	}
	userIDRes, err := ValidateJWT(jwt, tokenSecret)
	if err != nil {
		t.Error(err)
	}
	if !(userID == userIDRes) {
		t.Errorf(`User ID did not match ValidateJWT result`)
	}
}
func TestJWTExpired(t *testing.T) {
	userID := uuid.New()
	tokenSecret := "0rgasbord-org4asbord"
	jwt, err := MakeJWT(userID, tokenSecret)
	if err != nil {
		t.Error(err)
	}
	time.Sleep(500 * time.Millisecond)
	_, err = ValidateJWT(jwt, tokenSecret)
	if err == nil {
		t.Errorf(`Validate should have failed due to expired token:\n  %v`, err)
	}
}
func TestJWTWrongSecret(t *testing.T) {
	userID := uuid.New()
	tokenSecret := "0rgasbord-org4asbord"
	jwt, err := MakeJWT(userID, tokenSecret)
	if err != nil {
		t.Error(err)
	}
	_, err = ValidateJWT(jwt, "r4tly-fe4st")
	if err == nil {
		t.Errorf(`Validate should have failed due to incorrect tokenSecret:\n  %v`, err)
	}
}
func TestGetBearerToken(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	secret := "cook1e-Crumb$"
	req.Header.Set("Authorization", secret)
	returnAuth, err := GetBearerToken(req.Header)
	if (err != nil) || (returnAuth != secret) {
		t.Errorf(`TestGetBearerToken failed: %v`, err)
	}
}
func TestGetBearerTokenMismatch(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	secret := "cook1e-Crumb$"
	req.Header.Set("Authorization", secret+"$")
	returnAuth, err := GetBearerToken(req.Header)
	if (err != nil) || (returnAuth == secret) {
		t.Errorf(`TestGetBearerToken failed: %v`, err)
	}
}
