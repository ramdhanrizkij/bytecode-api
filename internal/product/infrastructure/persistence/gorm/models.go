package gorm

import "time"

type CategoryModel struct {
	ID   string `gorm:"column:id;type:uuid;primaryKey"`
	Name string `gorm:"column:name"`
}

func (CategoryModel) TableName() string {
	return "categories"
}

type ProductModel struct {
	ID          string        `gorm:"column:id;type:uuid;primaryKey"`
	CategoryID  string        `gorm:"column:category_id;type:uuid"`
	Category    CategoryModel `gorm:"foreignKey:CategoryID;references:ID"`
	Name        string        `gorm:"column:name"`
	Slug        string        `gorm:"column:slug"`
	Description *string       `gorm:"column:description"`
	SKU         string        `gorm:"column:sku"`
	Price       int64         `gorm:"column:price"`
	Stock       int           `gorm:"column:stock"`
	IsActive    bool          `gorm:"column:is_active"`
	CreatedAt   time.Time     `gorm:"column:created_at"`
	UpdatedAt   time.Time     `gorm:"column:updated_at"`
}

func (ProductModel) TableName() string {
	return "products"
}
