package repository

import (
	"context"

	"github.com/google/uuid"

	"github.com/ramdhanrizki/bytecode-api/internal/identity/domain/entity"
)

type RoleRepository interface {
	Create(ctx context.Context, role *entity.Role) error
	Update(ctx context.Context, role *entity.Role) error
	Delete(ctx context.Context, id uuid.UUID) error
	FindByID(ctx context.Context, id uuid.UUID) (*entity.Role, error)
	FindByIDWithPermissions(ctx context.Context, id uuid.UUID) (*entity.Role, error)
	FindByIDs(ctx context.Context, ids []uuid.UUID) ([]entity.Role, error)
	FindByName(ctx context.Context, name string) (*entity.Role, error)
	List(ctx context.Context, options ListOptions) ([]entity.Role, int, error)
}
