package domain

// CreateUserRequest is the payload for creating a new user.
type CreateUserRequest struct {
	Name           string  `json:"name"            validate:"required,min=2,max=100" example:"Jane Doe"`
	Email          string  `json:"email"           validate:"required,email" example:"jane@example.com"`
	Password       string  `json:"password"        validate:"required,min=8,max=50" example:"secret123"`
	RoleID         string  `json:"role_id"         validate:"required,uuid" example:"018f7606-a3f7-7c40-8e4b-2d47c6e04c8d"`
	ProfilePicture *string `json:"profile_picture" validate:"omitempty,max=500" example:"profiles/jane.png"`
	IsActive       *bool   `json:"is_active" example:"true"` // Use pointer to allow false-value validation
}

// UpdateUserRequest is the payload for updating an existing user.
type UpdateUserRequest struct {
	Name           string  `json:"name"            validate:"required,min=2,max=100" example:"Jane Doe"`
	Email          string  `json:"email"           validate:"required,email" example:"jane@example.com"`
	Password       string  `json:"password"        validate:"omitempty,min=8,max=50" example:"secret123"` // Optional on update
	RoleID         string  `json:"role_id"         validate:"required,uuid" example:"018f7606-a3f7-7c40-8e4b-2d47c6e04c8d"`
	ProfilePicture *string `json:"profile_picture" validate:"omitempty,max=500" example:"profiles/jane.png"`
	IsActive       *bool   `json:"is_active" example:"true"`
}

// UpdateProfileRequest is the payload for users updating their own profile.
type UpdateProfileRequest struct {
	Name           string  `json:"name"            validate:"required,min=2,max=100" example:"Jane Doe"`
	Email          string  `json:"email"           validate:"required,email" example:"jane@example.com"`
	ProfilePicture *string `json:"profile_picture" validate:"omitempty,max=500" example:"profiles/jane.png"`
}

// UserDetailResponse is the data returned for user details.
type UserDetailResponse struct {
	ID             string                  `json:"id" example:"018f7606-a3f7-7c40-8e4b-2d47c6e04c8d"`
	Name           string                  `json:"name" example:"Jane Doe"`
	Email          string                  `json:"email" example:"jane@example.com"`
	ProfilePicture *ProfilePictureResponse `json:"profile_picture,omitempty"`
	IsActive       bool                    `json:"is_active" example:"true"`
	Role           RoleInfo                `json:"role"`
	CreatedAt      string                  `json:"created_at" example:"2026-05-13T10:00:00+07:00"`
	UpdatedAt      string                  `json:"updated_at" example:"2026-05-13T10:00:00+07:00"`
}

// ProfilePictureResponse describes the resolved profile picture asset.
type ProfilePictureResponse struct {
	Bucket string `json:"bucket" example:"uploads"`
	Key    string `json:"key" example:"profiles/jane.png"`
	URL    string `json:"url" example:"http://localhost:8080/storage/profiles/jane.png"`
}

// RoleInfo provides basic role details in user responses.
type RoleInfo struct {
	ID   string `json:"id" example:"018f7606-a3f7-7c40-8e4b-2d47c6e04c8d"`
	Name string `json:"name" example:"user"`
}
