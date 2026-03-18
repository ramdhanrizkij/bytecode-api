package middleware

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	identityDomain "github.com/ramdhanrizki/bytecode-api/internal/identity/domain"
	"github.com/ramdhanrizki/bytecode-api/internal/identity/domain/repository"
	sharedAuth "github.com/ramdhanrizki/bytecode-api/internal/shared/auth"
	sharedResponse "github.com/ramdhanrizki/bytecode-api/internal/shared/response"
)

type PermissionMiddleware struct {
	users repository.UserRepository
}

func NewPermissionMiddleware(users repository.UserRepository) *PermissionMiddleware {
	return &PermissionMiddleware{users: users}
}

func (m *PermissionMiddleware) RequirePermission(permission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authUser, ok := sharedAuth.GetAuthenticatedUser(c)
		if !ok {
			sharedResponse.Error(c, identityDomain.ErrUnauthenticated)
			c.Abort()
			return
		}

		permissions, roles, err := m.permissionsForUser(c.Request.Context(), authUser.ID)
		if err != nil {
			sharedResponse.Error(c, err)
			c.Abort()
			return
		}

		authUser.Roles = roles
		authUser.Permissions = permissions
		sharedAuth.SetAuthenticatedUser(c, authUser)

		if !hasPermission(permissions, permission) {
			sharedResponse.Error(c, identityDomain.ErrForbidden)
			c.Abort()
			return
		}

		c.Next()
	}
}

func (m *PermissionMiddleware) permissionsForUser(ctx context.Context, userID uuid.UUID) ([]string, []string, error) {
	user, err := m.users.FindByIDWithRolesAndPermissions(ctx, userID)
	if err != nil {
		return nil, nil, err
	}
	if user == nil {
		return nil, nil, identityDomain.ErrUnauthenticated
	}

	permissionSet := make(map[string]struct{})
	permissions := make([]string, 0)
	roles := make([]string, 0, len(user.Roles))
	for _, role := range user.Roles {
		roles = append(roles, role.Name)
		for _, permission := range role.Permissions {
			if _, exists := permissionSet[permission.Name]; exists {
				continue
			}
			permissionSet[permission.Name] = struct{}{}
			permissions = append(permissions, permission.Name)
		}
	}

	return permissions, roles, nil
}

func hasPermission(permissions []string, permission string) bool {
	for _, candidate := range permissions {
		if candidate == permission {
			return true
		}
	}
	return false
}
