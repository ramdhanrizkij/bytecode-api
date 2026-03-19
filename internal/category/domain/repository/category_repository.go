package repository

import (
	"context"

	"github.com/google/uuid"

	"github.com/ramdhanrizki/bytecode-api/internal/category/domain/entity"
)

type ListOptions struct {
	Page     int
	Limit    int
	Search   string
	Sort     string
	Order    string
	IsActive *bool
}

type CategoryRepository interface {
	Create(ctx context.Context, category *entity.Category) error
	Update(ctx context.Context, category *entity.Category) error
	Delete(ctx context.Context, id uuid.UUID) error
	FindByID(ctx context.Context, id uuid.UUID) (*entity.Category, error)
	FindByName(ctx context.Context, name string) (*entity.Category, error)
	FindBySlug(ctx context.Context, slug string) (*entity.Category, error)
	List(ctx context.Context, options ListOptions) ([]entity.Category, int, error)
}
