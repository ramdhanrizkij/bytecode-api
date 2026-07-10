package model

import (
	"time"

	"github.com/google/uuid"
)

// User represents an application user with role-based access.
type User struct {
	ID             uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id" swaggertype:"string" format:"uuid" example:"018f7606-a3f7-7c40-8e4b-2d47c6e04c8d"`
	Name           string     `gorm:"type:varchar(100);not null" json:"name" example:"Jane Doe"`
	Email          string     `gorm:"type:varchar(100);uniqueIndex;not null" json:"email" example:"jane@example.com"`
	Password       string     `gorm:"type:varchar(255);not null" json:"-"`
	ProfilePicture *string    `gorm:"type:text" json:"profile_picture" example:"profiles/jane.png"`
	RoleID         *uuid.UUID `gorm:"type:uuid" json:"role_id" swaggertype:"string" format:"uuid" example:"018f7606-a3f7-7c40-8e4b-2d47c6e04c8d"`
	Role           *Role      `gorm:"foreignKey:RoleID" json:"role,omitempty"`
	IsActive       bool       `gorm:"default:true" json:"is_active" example:"true"`
	CreatedAt      time.Time  `json:"created_at" example:"2026-05-13T10:00:00+07:00"`
	UpdatedAt      time.Time  `json:"updated_at" example:"2026-05-13T10:00:00+07:00"`
}

// TableName returns the table name for the User model.
func (User) TableName() string {
	return "users"
}
