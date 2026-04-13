package model

import (
	"time"

	"github.com/google/uuid"
)

// User represents an application user with role-based access.
type User struct {
	ID        uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Name      string     `gorm:"type:varchar(100);not null" json:"name"`
	Email     string     `gorm:"type:varchar(100);uniqueIndex;not null" json:"email"`
	Password  string     `gorm:"type:varchar(255);not null" json:"-"`
	RoleID    *uuid.UUID `gorm:"type:uuid" json:"role_id"`
	Role      *Role      `gorm:"foreignKey:RoleID" json:"role,omitempty"`
	IsActive  bool       `gorm:"default:true" json:"is_active"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

// TableName returns the table name for the User model.
func (User) TableName() string {
	return "users"
}
