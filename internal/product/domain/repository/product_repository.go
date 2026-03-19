package repository

import (
	"context"

	"github.com/google/uuid"

	"github.com/ramdhanrizki/bytecode-api/internal/product/domain/entity"
)

type ListOptions struct {
	Page       int
	Limit      int
	Search     string
	Sort       string
	Order      string
	CategoryID *uuid.UUID
	IsActive   *bool
}

type ProductRepository interface {
	Create(ctx context.Context, product *entity.Product) error
	Update(ctx context.Context, product *entity.Product) error
	Delete(ctx context.Context, id uuid.UUID) error
	FindByID(ctx context.Context, id uuid.UUID) (*entity.Product, error)
	FindBySKU(ctx context.Context, sku string) (*entity.Product, error)
	FindBySlug(ctx context.Context, slug string) (*entity.Product, error)
	List(ctx context.Context, options ListOptions) ([]entity.Product, int, error)
}
