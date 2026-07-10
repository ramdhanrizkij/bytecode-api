package model

import (
	"time"

	"github.com/google/uuid"
)

// RolePermission is the explicit join table for the many-to-many relationship
// between roles and permissions.
type RolePermission struct {
	ID           uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id" swaggertype:"string" format:"uuid" example:"018f7606-a3f7-7c40-8e4b-2d47c6e04c8d"`
	RoleID       uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_role_permission" json:"role_id" swaggertype:"string" format:"uuid" example:"018f7606-a3f7-7c40-8e4b-2d47c6e04c8d"`
	PermissionID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_role_permission" json:"permission_id" swaggertype:"string" format:"uuid" example:"018f7606-a3f7-7c40-8e4b-2d47c6e04c8d"`
	CreatedAt    time.Time `json:"created_at" example:"2026-05-13T10:00:00+07:00"`
}

// TableName returns the table name for the RolePermission model.
func (RolePermission) TableName() string {
	return "role_permissions"
}
