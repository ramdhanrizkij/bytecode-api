package domain

// RegisterRequest is the payload for the POST /auth/register endpoint.
// Pure Go struct — no Fiber or any HTTP framework import allowed.
type RegisterRequest struct {
	Name     string `json:"name"     validate:"required,min=2,max=100"`
	Email    string `json:"email"    validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=50"`
}

// LoginRequest is the payload for the POST /auth/login endpoint.
type LoginRequest struct {
	Email    string `json:"email"    validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// AuthResponse is returned after a successful register or login.
type AuthResponse struct {
	User  UserResponse `json:"user"`
	Token string       `json:"token"`
}

// UserResponse is the user data embedded inside AuthResponse.
type UserResponse struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	RoleName string `json:"role_name"`
	IsActive bool   `json:"is_active"`
}
