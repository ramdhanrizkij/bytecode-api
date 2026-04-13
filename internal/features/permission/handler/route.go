package handler

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"github.com/ramdhanrizkij/bytecode-api/internal/core/middleware"
)

// RegisterRoutes registers all permission-related routes on the provided router group.
func RegisterRoutes(router fiber.Router, handler *PermissionHTTPHandler, db *gorm.DB, jwtSecret string) {
	perms := router.Group("/permissions", middleware.JWTAuth(jwtSecret))

	perms.Get("/", middleware.RequirePermission(db, "permissions.view"), handler.GetAll)
	perms.Get("/:id", middleware.RequirePermission(db, "permissions.view"), handler.GetByID)
	perms.Post("/", middleware.RequirePermission(db, "permissions.create"), handler.Create)
	perms.Put("/:id", middleware.RequirePermission(db, "permissions.edit"), handler.Update)
	perms.Delete("/:id", middleware.RequirePermission(db, "permissions.delete"), handler.Delete)
}
