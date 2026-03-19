package http

import (
	"github.com/gin-gonic/gin"

	"github.com/ramdhanrizki/bytecode-api/internal/category/delivery/http/handler"
	identityMiddleware "github.com/ramdhanrizki/bytecode-api/internal/identity/delivery/http/middleware"
)

type RouterDependencies struct {
	CategoryHandler      *handler.CategoryHandler
	AuthMiddleware       *identityMiddleware.AuthMiddleware
	PermissionMiddleware *identityMiddleware.PermissionMiddleware
}

func RegisterRoutes(router *gin.Engine, deps RouterDependencies) {
	api := router.Group("/api/v1")

	categoryRoutes := api.Group("/categories")
	categoryRoutes.Use(deps.AuthMiddleware.RequireAuth())
	categoryRoutes.GET("", deps.PermissionMiddleware.RequirePermission("categories.read"), deps.CategoryHandler.List)
	categoryRoutes.GET("/:id", deps.PermissionMiddleware.RequirePermission("categories.read"), deps.CategoryHandler.Get)

	adminRoutes := api.Group("/admin/categories")
	adminRoutes.Use(deps.AuthMiddleware.RequireAuth())
	adminRoutes.GET("", deps.PermissionMiddleware.RequirePermission("categories.read"), deps.CategoryHandler.List)
	adminRoutes.POST("", deps.PermissionMiddleware.RequirePermission("categories.create"), deps.CategoryHandler.Create)
	adminRoutes.GET("/:id", deps.PermissionMiddleware.RequirePermission("categories.read"), deps.CategoryHandler.Get)
	adminRoutes.PUT("/:id", deps.PermissionMiddleware.RequirePermission("categories.update"), deps.CategoryHandler.Update)
	adminRoutes.DELETE("/:id", deps.PermissionMiddleware.RequirePermission("categories.delete"), deps.CategoryHandler.Delete)
}
