package repository

import (
	"context"

	"github.com/google/uuid"
)

type RolePermissionRepository interface {
	ReplacePermissions(ctx context.Context, roleID uuid.UUID, permissionIDs []uuid.UUID) error
}
