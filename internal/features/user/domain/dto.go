package domain

// CreateUserRequest is the payload for creating a new user.
type CreateUserRequest struct {
	Name     string `json:"name"     validate:"required,min=2,max=100"`
	Email    string `json:"email"    validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=50"`
	RoleID   string `json:"role_id"  validate:"required,uuid"`
	IsActive *bool  `json:"is_active"` // Use pointer to allow false-value validation
}

// UpdateUserRequest is the payload for updating an existing user.
type UpdateUserRequest struct {
	Name     string `json:"name"     validate:"required,min=2,max=100"`
	Email    string `json:"email"    validate:"required,email"`
	Password string `json:"password" validate:"omitempty,min=8,max=50"` // Optional on update
	RoleID   string `json:"role_id"  validate:"required,uuid"`
	IsActive *bool  `json:"is_active"`
}

// UserDetailResponse is the data returned for user details.
type UserDetailResponse struct {
	ID        string   `json:"id"`
	Name      string   `json:"name"`
	Email     string   `json:"email"`
	IsActive  bool     `json:"is_active"`
	Role      RoleInfo `json:"role"`
	CreatedAt string   `json:"created_at"`
	UpdatedAt string   `json:"updated_at"`
}

// RoleInfo provides basic role details in user responses.
type RoleInfo struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
