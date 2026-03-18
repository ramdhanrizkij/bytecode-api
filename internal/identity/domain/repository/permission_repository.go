package repository

import (
	"context"

	"github.com/google/uuid"

	"github.com/ramdhanrizki/bytecode-api/internal/identity/domain/entity"
)

type PermissionRepository interface {
	Create(ctx context.Context, permission *entity.Permission) error
	Update(ctx context.Context, permission *entity.Permission) error
	Delete(ctx context.Context, id uuid.UUID) error
	FindByID(ctx context.Context, id uuid.UUID) (*entity.Permission, error)
	FindByIDs(ctx context.Context, ids []uuid.UUID) ([]entity.Permission, error)
	FindByName(ctx context.Context, name string) (*entity.Permission, error)
	List(ctx context.Context, options ListOptions) ([]entity.Permission, int, error)
}
