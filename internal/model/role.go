package model

import (
	"time"

	"github.com/google/uuid"
)

// Role represents a user role (e.g. superadmin, admin, user).
type Role struct {
	ID          uuid.UUID    `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id" swaggertype:"string" format:"uuid" example:"018f7606-a3f7-7c40-8e4b-2d47c6e04c8d"`
	Name        string       `gorm:"type:varchar(50);uniqueIndex;not null" json:"name" example:"manager"`
	Description string       `gorm:"type:text" json:"description" example:"Can manage operational resources"`
	GuardName   string       `gorm:"type:varchar(50);default:'api'" json:"guard_name" example:"api"`
	Permissions []Permission `gorm:"many2many:role_permissions;" json:"permissions,omitempty"`
	CreatedAt   time.Time    `json:"created_at" example:"2026-05-13T10:00:00+07:00"`
	UpdatedAt   time.Time    `json:"updated_at" example:"2026-05-13T10:00:00+07:00"`
}

// TableName returns the table name for the Role model.
func (Role) TableName() string {
	return "roles"
}
