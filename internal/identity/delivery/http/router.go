package http

import (
	"github.com/gin-gonic/gin"

	"github.com/ramdhanrizki/bytecode-api/internal/identity/delivery/http/handler"
	"github.com/ramdhanrizki/bytecode-api/internal/identity/delivery/http/middleware"
)

type RouterDependencies struct {
	AuthHandler          *handler.AuthHandler
	ProfileHandler       *handler.ProfileHandler
	AdminUserHandler     *handler.AdminUserHandler
	RoleHandler          *handler.RoleHandler
	PermissionHandler    *handler.PermissionHandler
	AuthMiddleware       *middleware.AuthMiddleware
	PermissionMiddleware *middleware.PermissionMiddleware
}

func RegisterRoutes(router *gin.Engine, deps RouterDependencies) {
	api := router.Group("/api/v1")

	authRoutes := api.Group("/auth")
	authRoutes.POST("/register", deps.AuthHandler.Register)
	authRoutes.POST("/verify-email", deps.AuthHandler.VerifyEmail)
	authRoutes.POST("/login", deps.AuthHandler.Login)
	authRoutes.POST("/refresh", deps.AuthHandler.Refresh)

	profileRoutes := api.Group("/profile")
	profileRoutes.Use(deps.AuthMiddleware.RequireAuth())
	profileRoutes.GET("", deps.PermissionMiddleware.RequirePermission("profile.read"), deps.ProfileHandler.GetCurrent)
	profileRoutes.PUT("", deps.PermissionMiddleware.RequirePermission("profile.update"), deps.ProfileHandler.UpdateCurrent)

	adminRoutes := api.Group("/admin")
	adminRoutes.Use(deps.AuthMiddleware.RequireAuth())

	adminRoutes.GET("/users", deps.PermissionMiddleware.RequirePermission("users.read"), deps.AdminUserHandler.List)
	adminRoutes.POST("/users", deps.PermissionMiddleware.RequirePermission("users.create"), deps.AdminUserHandler.Create)
	adminRoutes.GET("/users/:id", deps.PermissionMiddleware.RequirePermission("users.read"), deps.AdminUserHandler.Get)
	adminRoutes.PUT("/users/:id", deps.PermissionMiddleware.RequirePermission("users.update"), deps.AdminUserHandler.Update)
	adminRoutes.DELETE("/users/:id", deps.PermissionMiddleware.RequirePermission("users.delete"), deps.AdminUserHandler.Delete)
	adminRoutes.PUT("/users/:id/roles", deps.PermissionMiddleware.RequirePermission("users.update"), deps.AdminUserHandler.AssignRoles)

	adminRoutes.GET("/roles", deps.PermissionMiddleware.RequirePermission("roles.read"), deps.RoleHandler.List)
	adminRoutes.POST("/roles", deps.PermissionMiddleware.RequirePermission("roles.create"), deps.RoleHandler.Create)
	adminRoutes.GET("/roles/:id", deps.PermissionMiddleware.RequirePermission("roles.read"), deps.RoleHandler.Get)
	adminRoutes.PUT("/roles/:id", deps.PermissionMiddleware.RequirePermission("roles.update"), deps.RoleHandler.Update)
	adminRoutes.DELETE("/roles/:id", deps.PermissionMiddleware.RequirePermission("roles.delete"), deps.RoleHandler.Delete)
	adminRoutes.PUT("/roles/:id/permissions", deps.PermissionMiddleware.RequirePermission("roles.update"), deps.RoleHandler.AssignPermissions)

	adminRoutes.GET("/permissions", deps.PermissionMiddleware.RequirePermission("permissions.read"), deps.PermissionHandler.List)
	adminRoutes.POST("/permissions", deps.PermissionMiddleware.RequirePermission("permissions.create"), deps.PermissionHandler.Create)
	adminRoutes.GET("/permissions/:id", deps.PermissionMiddleware.RequirePermission("permissions.read"), deps.PermissionHandler.Get)
	adminRoutes.PUT("/permissions/:id", deps.PermissionMiddleware.RequirePermission("permissions.update"), deps.PermissionHandler.Update)
	adminRoutes.DELETE("/permissions/:id", deps.PermissionMiddleware.RequirePermission("permissions.delete"), deps.PermissionHandler.Delete)
}
