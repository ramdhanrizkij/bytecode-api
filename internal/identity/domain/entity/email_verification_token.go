package entity

import (
	"time"

	"github.com/google/uuid"
)

type EmailVerificationToken struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	Token     string
	ExpiresAt time.Time
	UsedAt    *time.Time
	CreatedAt time.Time
}

func (t EmailVerificationToken) IsExpired(now time.Time) bool {
	return now.After(t.ExpiresAt)
}

func (t EmailVerificationToken) IsUsed() bool {
	return t.UsedAt != nil
}
