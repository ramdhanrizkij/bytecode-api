package dto

import "time"

type RegisterInput struct {
	FullName string
	Email    string
	Password string
}

type RegisterOutput struct {
	User                 UserSummary
	VerificationQueued   bool
	VerificationTokenTTL time.Time
}

type VerifyEmailInput struct {
	Token string
}

type LoginInput struct {
	Email    string
	Password string
}

type RefreshTokenInput struct {
	RefreshToken string
}

type AuthTokens struct {
	AccessToken           string
	RefreshToken          string
	TokenType             string
	AccessTokenExpiresAt  time.Time
	RefreshTokenExpiresAt time.Time
}

type AuthOutput struct {
	User   UserSummary
	Tokens AuthTokens
}
