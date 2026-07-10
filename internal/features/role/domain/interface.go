package domain

import (
	"context"

	"github.com/ramdhanrizkij/bytecode-api/internal/model"
	"github.com/ramdhanrizkij/bytecode-api/internal/shared/pagination"
	"github.com/ramdhanrizkij/bytecode-api/internal/shared/response"
)

// RoleRepository defines the data-access contract for the role feature.
type RoleRepository interface {
	FindAll(ctx context.Context, pq *pagination.PaginationQuery) ([]model.Role, int64, error)
	FindByID(ctx context.Context, id string) (*model.Role, error)
	FindByName(ctx context.Context, name string) (*model.Role, error)
	Create(ctx context.Context, role *model.Role) error
	Update(ctx context.Context, role *model.Role) error
	Delete(ctx context.Context, id string) error
	AssignPermissions(ctx context.Context, roleID string, permissionIDs []string) error
	RemovePermissions(ctx context.Context, roleID string, permissionIDs []string) error
}

// RoleService defines the business-logic contract for the role feature.
type RoleService interface {
	GetAll(ctx context.Context, pq *pagination.PaginationQuery) ([]RoleResponse, *response.PaginationMeta, error)
	GetByID(ctx context.Context, id string) (*RoleResponse, error)
	Create(ctx context.Context, req *CreateRoleRequest) (*RoleResponse, error)
	Update(ctx context.Context, id string, req *UpdateRoleRequest) (*RoleResponse, error)
	Delete(ctx context.Context, id string) error
	AssignPermissions(ctx context.Context, id string, req *AssignPermissionsRequest) error
	RemovePermissions(ctx context.Context, id string, req *AssignPermissionsRequest) error
}
