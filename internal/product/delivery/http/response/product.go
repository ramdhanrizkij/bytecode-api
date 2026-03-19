package response

import (
	"time"

	"github.com/ramdhanrizki/bytecode-api/internal/product/application/dto"
)

type ProductResponse struct {
	ID           string    `json:"id"`
	CategoryID   string    `json:"category_id"`
	CategoryName string    `json:"category_name"`
	Name         string    `json:"name"`
	Slug         string    `json:"slug"`
	Description  *string   `json:"description,omitempty"`
	SKU          string    `json:"sku"`
	Price        int64     `json:"price"`
	Stock        int       `json:"stock"`
	IsActive     bool      `json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func FromProductSummary(product dto.ProductSummary) ProductResponse {
	return ProductResponse{
		ID:           product.ID,
		CategoryID:   product.CategoryID,
		CategoryName: product.CategoryName,
		Name:         product.Name,
		Slug:         product.Slug,
		Description:  product.Description,
		SKU:          product.SKU,
		Price:        product.Price,
		Stock:        product.Stock,
		IsActive:     product.IsActive,
		CreatedAt:    product.CreatedAt,
		UpdatedAt:    product.UpdatedAt,
	}
}

func FromProductSummaries(products []dto.ProductSummary) []ProductResponse {
	items := make([]ProductResponse, 0, len(products))
	for _, product := range products {
		items = append(items, FromProductSummary(product))
	}
	return items
}
