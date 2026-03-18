package response

import (
	"time"

	"github.com/ramdhanrizki/bytecode-api/internal/identity/application/dto"
)

type RegisterResponse struct {
	User                  UserResponse `json:"user"`
	VerificationQueued    bool         `json:"verification_queued"`
	VerificationExpiresAt time.Time    `json:"verification_expires_at"`
}

type AuthResponse struct {
	User   UserResponse   `json:"user"`
	Tokens TokensResponse `json:"tokens"`
}

type TokensResponse struct {
	AccessToken           string    `json:"access_token"`
	RefreshToken          string    `json:"refresh_token"`
	TokenType             string    `json:"token_type"`
	AccessTokenExpiresAt  time.Time `json:"access_token_expires_at"`
	RefreshTokenExpiresAt time.Time `json:"refresh_token_expires_at"`
}

type UserResponse struct {
	ID              string   `json:"id"`
	FullName        string   `json:"full_name"`
	Email           string   `json:"email"`
	IsEmailVerified bool     `json:"is_email_verified"`
	IsActive        bool     `json:"is_active"`
	Roles           []string `json:"roles"`
}

func FromRegisterOutput(output dto.RegisterOutput) RegisterResponse {
	return RegisterResponse{
		User:                  FromUserSummary(output.User),
		VerificationQueued:    output.VerificationQueued,
		VerificationExpiresAt: output.VerificationTokenTTL,
	}
}

func FromAuthOutput(output dto.AuthOutput) AuthResponse {
	return AuthResponse{
		User: FromUserSummary(output.User),
		Tokens: TokensResponse{
			AccessToken:           output.Tokens.AccessToken,
			RefreshToken:          output.Tokens.RefreshToken,
			TokenType:             output.Tokens.TokenType,
			AccessTokenExpiresAt:  output.Tokens.AccessTokenExpiresAt,
			RefreshTokenExpiresAt: output.Tokens.RefreshTokenExpiresAt,
		},
	}
}

func FromUserSummary(user dto.UserSummary) UserResponse {
	return UserResponse{
		ID:              user.ID,
		FullName:        user.FullName,
		Email:           user.Email,
		IsEmailVerified: user.IsEmailVerified,
		IsActive:        user.IsActive,
		Roles:           user.Roles,
	}
}
