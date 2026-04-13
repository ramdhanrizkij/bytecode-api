package middleware

import (
	"strings"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"github.com/ramdhanrizkij/bytecode-api/internal/model"
	"github.com/ramdhanrizkij/bytecode-api/internal/shared/response"
)

// permissionCache stores permissions for a role name.
// Key: roleName, Value: cachedPermissions
var (
	permCache sync.Map
	cacheTTL  = 5 * time.Minute
)

type cachedPermissions struct {
	permissions []string
	expiry      time.Time
}

// RequireRole ensures the authenticated user has one of the specified roles.
func RequireRole(roles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		claims := GetCurrentUser(c)
		if claims == nil {
			return response.Error(c, fiber.StatusUnauthorized, "unauthorized")
		}

		for _, role := range roles {
			if strings.EqualFold(claims.RoleName, role) {
				return c.Next()
			}
		}

		return response.Error(c, fiber.StatusForbidden, "insufficient role permissions")
	}
}

// RequirePermission ensures the authenticated user's role has one of the specified permissions.
func RequirePermission(db *gorm.DB, permissions ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		claims := GetCurrentUser(c)
		if claims == nil {
			return response.Error(c, fiber.StatusUnauthorized, "unauthorized")
		}

		// Superadmin bypass
		if strings.EqualFold(claims.RoleName, "superadmin") {
			return c.Next()
		}

		rolePerms, err := getRolePermissions(db, claims.RoleName)
		if err != nil {
			return response.Error(c, fiber.StatusInternalServerError, "failed to verify permissions")
		}

		permMap := make(map[string]bool)
		for _, p := range rolePerms {
			permMap[p] = true
		}

		for _, required := range permissions {
			if permMap[required] {
				return c.Next()
			}
		}

		return response.Error(c, fiber.StatusForbidden, "insufficient permissions")
	}
}

func getRolePermissions(db *gorm.DB, roleName string) ([]string, error) {
	// Check cache
	if val, ok := permCache.Load(roleName); ok {
		cp := val.(cachedPermissions)
		if time.Now().Before(cp.expiry) {
			return cp.permissions, nil
		}
	}

	// Cache miss or expired - query DB
	var role model.Role
	err := db.Preload("Permissions").Where("name = ?", roleName).First(&role).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return []string{}, nil
		}
		return nil, err
	}

	permNames := make([]string, 0, len(role.Permissions))
	for _, p := range role.Permissions {
		permNames = append(permNames, p.Name)
	}

	// Update cache
	permCache.Store(roleName, cachedPermissions{
		permissions: permNames,
		expiry:      time.Now().Add(cacheTTL),
	})

	return permNames, nil
}
