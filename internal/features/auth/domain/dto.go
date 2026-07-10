package domain

// RegisterRequest is the payload for the POST /auth/register endpoint.
// Pure Go struct — no Fiber or any HTTP framework import allowed.
type RegisterRequest struct {
	Name           string  `json:"name"            validate:"required,min=2,max=100" example:"Jane Doe"`
	Email          string  `json:"email"           validate:"required,email" example:"jane@example.com"`
	Password       string  `json:"password"        validate:"required,min=8,max=50" example:"secret123"`
	ProfilePicture *string `json:"profile_picture" validate:"omitempty,max=500" example:"profiles/jane.png"`
}

// LoginRequest is the payload for the POST /auth/login endpoint.
type LoginRequest struct {
	Email    string `json:"email"    validate:"required,email" example:"jane@example.com"`
	Password string `json:"password" validate:"required" example:"secret123"`
}

// RefreshTokenRequest is the payload for the POST /auth/refresh endpoint.
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required" example:"refresh-token-value"`
}

// LogoutRequest is the payload for the POST /auth/logout endpoint.
type LogoutRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required" example:"refresh-token-value"`
}

// AuthResponse is returned after a successful register or login.
type AuthResponse struct {
	User         UserResponse `json:"user"`
	Token        string       `json:"token" example:"access-token-value"`
	RefreshToken string       `json:"refresh_token" example:"refresh-token-value"`
}

// TokenResponse is returned after a successful refresh.
type TokenResponse struct {
	Token        string `json:"token" example:"access-token-value"`
	RefreshToken string `json:"refresh_token" example:"refresh-token-value"`
}

// UserResponse is the user data embedded inside AuthResponse.
type UserResponse struct {
	ID       string `json:"id" example:"018f7606-a3f7-7c40-8e4b-2d47c6e04c8d"`
	Name     string `json:"name" example:"Jane Doe"`
	Email    string `json:"email" example:"jane@example.com"`
	RoleName string `json:"role_name" example:"user"`
	IsActive bool   `json:"is_active" example:"true"`
}
