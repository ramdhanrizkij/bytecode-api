package http

import (
	"github.com/gin-gonic/gin"

	identityMiddleware "github.com/ramdhanrizki/bytecode-api/internal/identity/delivery/http/middleware"
	"github.com/ramdhanrizki/bytecode-api/internal/product/delivery/http/handler"
)

type RouterDependencies struct {
	ProductHandler       *handler.ProductHandler
	AuthMiddleware       *identityMiddleware.AuthMiddleware
	PermissionMiddleware *identityMiddleware.PermissionMiddleware
}

func RegisterRoutes(router *gin.Engine, deps RouterDependencies) {
	api := router.Group("/api/v1")

	productRoutes := api.Group("/products")
	productRoutes.Use(deps.AuthMiddleware.RequireAuth())
	productRoutes.GET("", deps.PermissionMiddleware.RequirePermission("products.read"), deps.ProductHandler.List)
	productRoutes.GET("/:id", deps.PermissionMiddleware.RequirePermission("products.read"), deps.ProductHandler.Get)

	adminRoutes := api.Group("/admin/products")
	adminRoutes.Use(deps.AuthMiddleware.RequireAuth())
	adminRoutes.GET("", deps.PermissionMiddleware.RequirePermission("products.read"), deps.ProductHandler.List)
	adminRoutes.POST("", deps.PermissionMiddleware.RequirePermission("products.create"), deps.ProductHandler.Create)
	adminRoutes.GET("/:id", deps.PermissionMiddleware.RequirePermission("products.read"), deps.ProductHandler.Get)
	adminRoutes.PUT("/:id", deps.PermissionMiddleware.RequirePermission("products.update"), deps.ProductHandler.Update)
	adminRoutes.DELETE("/:id", deps.PermissionMiddleware.RequirePermission("products.delete"), deps.ProductHandler.Delete)
}
