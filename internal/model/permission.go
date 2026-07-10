package model

import (
	"time"

	"github.com/google/uuid"
)

// Permission represents a single action a role can perform (e.g. "users:read").
type Permission struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id" swaggertype:"string" format:"uuid" example:"018f7606-a3f7-7c40-8e4b-2d47c6e04c8d"`
	Name        string    `gorm:"type:varchar(100);uniqueIndex;not null" json:"name" example:"users.view"`
	Description string    `gorm:"type:text" json:"description" example:"Can view users"`
	GuardName   string    `gorm:"type:varchar(50);default:'api'" json:"guard_name" example:"api"`
	CreatedAt   time.Time `json:"created_at" example:"2026-05-13T10:00:00+07:00"`
	UpdatedAt   time.Time `json:"updated_at" example:"2026-05-13T10:00:00+07:00"`
}

// TableName returns the table name for the Permission model.
func (Permission) TableName() string {
	return "permissions"
}
