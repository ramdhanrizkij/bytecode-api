package category

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	categoryService "github.com/ramdhanrizki/bytecode-api/internal/category/application/service"
	categoryHTTP "github.com/ramdhanrizki/bytecode-api/internal/category/delivery/http"
	categoryHandler "github.com/ramdhanrizki/bytecode-api/internal/category/delivery/http/handler"
	categoryGorm "github.com/ramdhanrizki/bytecode-api/internal/category/infrastructure/persistence/gorm"
	identityMiddleware "github.com/ramdhanrizki/bytecode-api/internal/identity/delivery/http/middleware"
	sharedLogger "github.com/ramdhanrizki/bytecode-api/internal/shared/logger"
)

type Module struct {
	CategoryHandler *categoryHandler.CategoryHandler
}

type Dependencies struct {
	Logger sharedLogger.Logger
	DB     *gorm.DB
}

func NewModule(deps Dependencies) *Module {
	repos := categoryGorm.NewRepositories(deps.DB)
	service := categoryService.NewCategoryService(repos.Categories())

	return &Module{
		CategoryHandler: categoryHandler.NewCategoryHandler(service),
	}
}

func (m *Module) RegisterRoutes(router *gin.Engine, authMiddleware *identityMiddleware.AuthMiddleware, permissionMiddleware *identityMiddleware.PermissionMiddleware) {
	categoryHTTP.RegisterRoutes(router, categoryHTTP.RouterDependencies{
		CategoryHandler:      m.CategoryHandler,
		AuthMiddleware:       authMiddleware,
		PermissionMiddleware: permissionMiddleware,
	})
}
