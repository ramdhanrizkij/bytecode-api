package product

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	categoryGorm "github.com/ramdhanrizki/bytecode-api/internal/category/infrastructure/persistence/gorm"
	identityMiddleware "github.com/ramdhanrizki/bytecode-api/internal/identity/delivery/http/middleware"
	productService "github.com/ramdhanrizki/bytecode-api/internal/product/application/service"
	productHTTP "github.com/ramdhanrizki/bytecode-api/internal/product/delivery/http"
	productHandler "github.com/ramdhanrizki/bytecode-api/internal/product/delivery/http/handler"
	productGorm "github.com/ramdhanrizki/bytecode-api/internal/product/infrastructure/persistence/gorm"
	sharedLogger "github.com/ramdhanrizki/bytecode-api/internal/shared/logger"
)

type Module struct {
	ProductHandler *productHandler.ProductHandler
}

type Dependencies struct {
	Logger sharedLogger.Logger
	DB     *gorm.DB
}

func NewModule(deps Dependencies) *Module {
	productRepos := productGorm.NewRepositories(deps.DB)
	categoryRepos := categoryGorm.NewRepositories(deps.DB)
	service := productService.NewProductService(productRepos.Products(), categoryRepos.Categories())

	return &Module{
		ProductHandler: productHandler.NewProductHandler(service),
	}
}

func (m *Module) RegisterRoutes(router *gin.Engine, authMiddleware *identityMiddleware.AuthMiddleware, permissionMiddleware *identityMiddleware.PermissionMiddleware) {
	productHTTP.RegisterRoutes(router, productHTTP.RouterDependencies{
		ProductHandler:       m.ProductHandler,
		AuthMiddleware:       authMiddleware,
		PermissionMiddleware: permissionMiddleware,
	})
}
