package repository

import (
	"context"

	"github.com/google/uuid"
)

type UserRoleRepository interface {
	AssignRole(ctx context.Context, userID, roleID uuid.UUID) error
	ReplaceRoles(ctx context.Context, userID uuid.UUID, roleIDs []uuid.UUID) error
}
