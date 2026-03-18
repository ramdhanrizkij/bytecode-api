package entity

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID              uuid.UUID
	FullName        string
	Email           string
	PasswordHash    string
	IsEmailVerified bool
	IsActive        bool
	Roles           []Role
	CreatedAt       time.Time
	UpdatedAt       time.Time
}
