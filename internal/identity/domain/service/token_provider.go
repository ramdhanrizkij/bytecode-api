package service

import (
	"time"

	"github.com/google/uuid"
)

type AccessTokenClaims struct {
	UserID uuid.UUID
	Email  string
	Roles  []string
}

type AccessToken struct {
	Token     string
	ExpiresAt time.Time
}

type ParsedAccessToken struct {
	UserID uuid.UUID
	Email  string
	Roles  []string
}

type AccessTokenProvider interface {
	Generate(claims AccessTokenClaims) (AccessToken, error)
	Parse(token string) (*ParsedAccessToken, error)
}
