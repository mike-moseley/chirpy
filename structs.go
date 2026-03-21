package main

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID          uuid.UUID `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Email       string    `json:"email"`
	IsChirpyRed bool      `json:"is_chirpy_red"`
}

type loginResponse struct {
	ID           uuid.UUID `json:"id"`
	Email        string    `json:"email"`
	AccessToken  string    `json:"token"`
	RefreshToken string    `json:"refresh_token"`
	IsChirpyRed  bool      `json:"is_chirpy_red"`
}

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

type RefreshToken struct {
	Token     string       `json:"token"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
	UserID    uuid.UUID    `json:"user_id"`
	ExpiresAt time.Time    `json:"expires_at"`
	RevokedAt sql.NullTime `json:"revoked_at"`
}
