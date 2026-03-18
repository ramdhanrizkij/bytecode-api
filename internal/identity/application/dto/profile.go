package dto

type UserSummary struct {
	ID              string
	FullName        string
	Email           string
	IsEmailVerified bool
	IsActive        bool
	Roles           []string
}

type UpdateProfileInput struct {
	UserID   string
	FullName string
}
