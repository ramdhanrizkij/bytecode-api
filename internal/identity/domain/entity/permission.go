package entity

import (
	"time"

	"github.com/google/uuid"
)

type Permission struct {
	ID          uuid.UUID
	Name        string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
