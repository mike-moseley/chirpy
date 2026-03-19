package auth

import (
	"testing"

	"github.com/alexedwards/argon2id"
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
