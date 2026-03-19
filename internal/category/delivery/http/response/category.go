package response

import (
	"time"

	"github.com/ramdhanrizki/bytecode-api/internal/category/application/dto"
)

type CategoryResponse struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Slug        string    `json:"slug"`
	Description *string   `json:"description,omitempty"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func FromCategorySummary(category dto.CategorySummary) CategoryResponse {
	return CategoryResponse{
		ID:          category.ID,
		Name:        category.Name,
		Slug:        category.Slug,
		Description: category.Description,
		IsActive:    category.IsActive,
		CreatedAt:   category.CreatedAt,
		UpdatedAt:   category.UpdatedAt,
	}
}

func FromCategorySummaries(categories []dto.CategorySummary) []CategoryResponse {
	items := make([]CategoryResponse, 0, len(categories))
	for _, category := range categories {
		items = append(items, FromCategorySummary(category))
	}
	return items
}
