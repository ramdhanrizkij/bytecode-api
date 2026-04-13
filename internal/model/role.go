package model

import (
	"time"

	"github.com/google/uuid"
)

// Role represents a user role (e.g. superadmin, admin, user).
type Role struct {
	ID          uuid.UUID    `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Name        string       `gorm:"type:varchar(50);uniqueIndex;not null" json:"name"`
	Description string       `gorm:"type:text" json:"description"`
	GuardName   string       `gorm:"type:varchar(50);default:'api'" json:"guard_name"`
	Permissions []Permission `gorm:"many2many:role_permissions;" json:"permissions,omitempty"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
}

// TableName returns the table name for the Role model.
func (Role) TableName() string {
	return "roles"
}
