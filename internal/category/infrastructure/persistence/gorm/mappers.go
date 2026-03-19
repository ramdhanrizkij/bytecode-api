package gorm

import (
	"github.com/google/uuid"

	"github.com/ramdhanrizki/bytecode-api/internal/category/domain/entity"
)

func toCategoryModel(category entity.Category) *CategoryModel {
	return &CategoryModel{
		ID:          category.ID.String(),
		Name:        category.Name,
		Slug:        category.Slug,
		Description: category.Description,
		IsActive:    category.IsActive,
		CreatedAt:   category.CreatedAt,
		UpdatedAt:   category.UpdatedAt,
	}
}

func toCategoryEntity(model CategoryModel) (*entity.Category, error) {
	id, err := uuid.Parse(model.ID)
	if err != nil {
		return nil, err
	}

	return &entity.Category{
		ID:          id,
		Name:        model.Name,
		Slug:        model.Slug,
		Description: model.Description,
		IsActive:    model.IsActive,
		CreatedAt:   model.CreatedAt,
		UpdatedAt:   model.UpdatedAt,
	}, nil
}
