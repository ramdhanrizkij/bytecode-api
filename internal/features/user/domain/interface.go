package domain

import (
	"context"

	"github.com/ramdhanrizkij/bytecode-api/internal/model"
	"github.com/ramdhanrizkij/bytecode-api/internal/shared/pagination"
	"github.com/ramdhanrizkij/bytecode-api/internal/shared/response"
)

// UserRepository defines the data access contract for users.
type UserRepository interface {
	FindAll(ctx context.Context, pq *pagination.PaginationQuery) ([]model.User, int64, error)
	FindByID(ctx context.Context, id string) (*model.User, error)
	FindByEmail(ctx context.Context, email string) (*model.User, error)
	Create(ctx context.Context, user *model.User) error
	Update(ctx context.Context, user *model.User) error
	Delete(ctx context.Context, id string) error
	GetPermissions(ctx context.Context, userID string) ([]model.Permission, error)
}

// UserService defines the business logic contract for users.
type UserService interface {
	GetAll(ctx context.Context, pq *pagination.PaginationQuery) ([]UserDetailResponse, *response.PaginationMeta, error)
	GetByID(ctx context.Context, id string) (*UserDetailResponse, error)
	Create(ctx context.Context, req *CreateUserRequest) (*UserDetailResponse, error)
	Update(ctx context.Context, id string, req *UpdateUserRequest) (*UserDetailResponse, error)
	Delete(ctx context.Context, currentUserID string, targetID string) error
	GetPermissions(ctx context.Context, userID string) ([]string, error)
}
