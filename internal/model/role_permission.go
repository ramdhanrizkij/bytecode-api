package model

import (
	"time"

	"github.com/google/uuid"
)

// RolePermission is the explicit join table for the many-to-many relationship
// between roles and permissions.
type RolePermission struct {
	ID           uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	RoleID       uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_role_permission" json:"role_id"`
	PermissionID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_role_permission" json:"permission_id"`
	CreatedAt    time.Time `json:"created_at"`
}

// TableName returns the table name for the RolePermission model.
func (RolePermission) TableName() string {
	return "role_permissions"
}
