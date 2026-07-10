package domain

import (
	"context"

	"github.com/ramdhanrizkij/bytecode-api/internal/model"
	"github.com/ramdhanrizkij/bytecode-api/internal/shared/pagination"
	"github.com/ramdhanrizkij/bytecode-api/internal/shared/response"
)

// PermissionRepository defines the data-access contract for the permission feature.
type PermissionRepository interface {
	FindAll(ctx context.Context, pq *pagination.PaginationQuery) ([]model.Permission, int64, error)
	FindByID(ctx context.Context, id string) (*model.Permission, error)
	FindByName(ctx context.Context, name string) (*model.Permission, error)
	FindByIDs(ctx context.Context, ids []string) ([]model.Permission, error)
	Create(ctx context.Context, permission *model.Permission) error
	Update(ctx context.Context, permission *model.Permission) error
	Delete(ctx context.Context, id string) error
}

// PermissionService defines the business-logic contract for the permission feature.
type PermissionService interface {
	GetAll(ctx context.Context, pq *pagination.PaginationQuery) ([]PermissionDetailResponse, *response.PaginationMeta, error)
	GetByID(ctx context.Context, id string) (*PermissionDetailResponse, error)
	Create(ctx context.Context, req *CreatePermissionRequest) (*PermissionDetailResponse, error)
	Update(ctx context.Context, id string, req *UpdatePermissionRequest) (*PermissionDetailResponse, error)
	Delete(ctx context.Context, id string) error
}
