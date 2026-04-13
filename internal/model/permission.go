package model

import (
	"time"

	"github.com/google/uuid"
)

// Permission represents a single action a role can perform (e.g. "users:read").
type Permission struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Name        string    `gorm:"type:varchar(100);uniqueIndex;not null" json:"name"`
	Description string    `gorm:"type:text" json:"description"`
	GuardName   string    `gorm:"type:varchar(50);default:'api'" json:"guard_name"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// TableName returns the table name for the Permission model.
func (Permission) TableName() string {
	return "permissions"
}
