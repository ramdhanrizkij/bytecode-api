package gorm

import (
	"github.com/google/uuid"

	"github.com/ramdhanrizki/bytecode-api/internal/product/domain/entity"
)

func toProductModel(product entity.Product) *ProductModel {
	return &ProductModel{
		ID:          product.ID.String(),
		CategoryID:  product.CategoryID.String(),
		Name:        product.Name,
		Slug:        product.Slug,
		Description: product.Description,
		SKU:         product.SKU,
		Price:       product.Price,
		Stock:       product.Stock,
		IsActive:    product.IsActive,
		CreatedAt:   product.CreatedAt,
		UpdatedAt:   product.UpdatedAt,
	}
}

func toProductEntity(model ProductModel) (*entity.Product, error) {
	id, err := uuid.Parse(model.ID)
	if err != nil {
		return nil, err
	}
	categoryID, err := uuid.Parse(model.CategoryID)
	if err != nil {
		return nil, err
	}

	return &entity.Product{
		ID:           id,
		CategoryID:   categoryID,
		CategoryName: model.Category.Name,
		Name:         model.Name,
		Slug:         model.Slug,
		Description:  model.Description,
		SKU:          model.SKU,
		Price:        model.Price,
		Stock:        model.Stock,
		IsActive:     model.IsActive,
		CreatedAt:    model.CreatedAt,
		UpdatedAt:    model.UpdatedAt,
	}, nil
}
