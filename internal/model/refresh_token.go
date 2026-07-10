package model

import (
	"time"

	"github.com/google/uuid"
)

// RefreshToken stores a hashed opaque refresh token for a user session.
type RefreshToken struct {
	ID        uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id" swaggertype:"string" format:"uuid" example:"018f7606-a3f7-7c40-8e4b-2d47c6e04c8d"`
	UserID    uuid.UUID  `gorm:"type:uuid;not null;index" json:"user_id" swaggertype:"string" format:"uuid" example:"018f7606-a3f7-7c40-8e4b-2d47c6e04c8d"`
	User      *User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
	TokenHash string     `gorm:"type:varchar(64);uniqueIndex;not null" json:"-"`
	ExpiresAt time.Time  `gorm:"not null;index" json:"expires_at" example:"2026-05-20T10:00:00+07:00"`
	RevokedAt *time.Time `json:"revoked_at,omitempty" example:"2026-05-13T10:00:00+07:00"`
	CreatedAt time.Time  `json:"created_at" example:"2026-05-13T10:00:00+07:00"`
	UpdatedAt time.Time  `json:"updated_at" example:"2026-05-13T10:00:00+07:00"`
}

// TableName returns the table name for the RefreshToken model.
func (RefreshToken) TableName() string {
	return "refresh_tokens"
}
