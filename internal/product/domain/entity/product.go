package entity

import (
	"time"

	"github.com/google/uuid"
)

type Product struct {
	ID           uuid.UUID
	CategoryID   uuid.UUID
	CategoryName string
	Name         string
	Slug         string
	Description  *string
	SKU          string
	Price        int64
	Stock        int
	IsActive     bool
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
