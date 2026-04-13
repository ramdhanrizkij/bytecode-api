package handler

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"github.com/ramdhanrizkij/bytecode-api/internal/core/middleware"
)

// RegisterRoutes registers all role-related routes on the provided router group.
func RegisterRoutes(router fiber.Router, handler *RoleHTTPHandler, db *gorm.DB, jwtSecret string) {
	roles := router.Group("/roles", middleware.JWTAuth(jwtSecret))

	roles.Get("/", middleware.RequirePermission(db, "roles.view"), handler.GetAll)
	roles.Get("/:id", middleware.RequirePermission(db, "roles.view"), handler.GetByID)
	roles.Post("/", middleware.RequirePermission(db, "roles.create"), handler.Create)
	roles.Put("/:id", middleware.RequirePermission(db, "roles.edit"), handler.Update)
	roles.Delete("/:id", middleware.RequirePermission(db, "roles.delete"), handler.Delete)
	roles.Post("/:id/permissions", middleware.RequirePermission(db, "roles.assign-permission"), handler.AssignPermissions)
	roles.Delete("/:id/permissions", middleware.RequirePermission(db, "roles.remove-permission"), handler.RemovePermissions)
}
