package gorm

import "time"

type CategoryModel struct {
	ID          string    `gorm:"column:id;type:uuid;primaryKey"`
	Name        string    `gorm:"column:name"`
	Slug        string    `gorm:"column:slug"`
	Description *string   `gorm:"column:description"`
	IsActive    bool      `gorm:"column:is_active"`
	CreatedAt   time.Time `gorm:"column:created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at"`
}

func (CategoryModel) TableName() string {
	return "categories"
}
