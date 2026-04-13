package handler

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"github.com/ramdhanrizkij/bytecode-api/internal/core/middleware"
)

// RegisterRoutes registers all user-related routes on the provided router.
func RegisterRoutes(router fiber.Router, handler *UserHTTPHandler, db *gorm.DB, jwtSecret string) {
	users := router.Group("/users", middleware.JWTAuth(jwtSecret))

	// Get current user profile - must be registered BEFORE /:id to avoid conflict.
	// These only require JWT authentication.
	users.Get("/me", handler.GetMe)
	users.Get("/me/permissions", handler.GetMePermissions)

	// Admin-level operations require specific permissions.
	users.Get("/", middleware.RequirePermission(db, "users.view"), handler.GetAll)
	users.Get("/:id", middleware.RequirePermission(db, "users.view"), handler.GetByID)
	users.Post("/", middleware.RequirePermission(db, "users.create"), handler.Create)
	users.Put("/:id", middleware.RequirePermission(db, "users.edit"), handler.Update)
	users.Delete("/:id", middleware.RequirePermission(db, "users.delete"), handler.Delete)
}
