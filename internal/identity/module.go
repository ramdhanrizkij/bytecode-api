package identity

import (
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/ramdhanrizki/bytecode-api/configs"
	appService "github.com/ramdhanrizki/bytecode-api/internal/identity/application/service"
	identityHTTP "github.com/ramdhanrizki/bytecode-api/internal/identity/delivery/http"
	identityHandler "github.com/ramdhanrizki/bytecode-api/internal/identity/delivery/http/handler"
	identityMiddleware "github.com/ramdhanrizki/bytecode-api/internal/identity/delivery/http/middleware"
	identityGorm "github.com/ramdhanrizki/bytecode-api/internal/identity/infrastructure/persistence/gorm"
	identityQueue "github.com/ramdhanrizki/bytecode-api/internal/identity/infrastructure/queue"
	platformCrypto "github.com/ramdhanrizki/bytecode-api/internal/platform/crypto"
	platformJWT "github.com/ramdhanrizki/bytecode-api/internal/platform/jwt"
	sharedLogger "github.com/ramdhanrizki/bytecode-api/internal/shared/logger"
	sharedQueue "github.com/ramdhanrizki/bytecode-api/internal/shared/queue"
)

type Module struct {
	AuthHandler          *identityHandler.AuthHandler
	ProfileHandler       *identityHandler.ProfileHandler
	AdminUserHandler     *identityHandler.AdminUserHandler
	RoleHandler          *identityHandler.RoleHandler
	PermissionHandler    *identityHandler.PermissionHandler
	AuthMiddleware       *identityMiddleware.AuthMiddleware
	PermissionMiddleware *identityMiddleware.PermissionMiddleware
}

type Dependencies struct {
	Config configs.Config
	Logger sharedLogger.Logger
	DB     *gorm.DB
	Queue  sharedQueue.Publisher
}

func NewModule(deps Dependencies) *Module {
	repos := identityGorm.NewRepositories(deps.DB)
	unitOfWork := identityGorm.NewUnitOfWork(deps.DB)
	passwordHasher := platformCrypto.NewBcryptHasher(0)
	tokenProvider := platformJWT.NewProvider(deps.Config.JWT)
	tokenGenerator := platformJWT.NewRandomTokenGenerator(32)
	jobPublisher := identityQueue.NewVerificationPublisher(deps.Queue)

	authService := appService.NewAuthService(appService.AuthServiceDependencies{
		Logger:              deps.Logger,
		Users:               repos.Users(),
		Roles:               repos.Roles(),
		RefreshTokens:       repos.RefreshTokens(),
		VerificationTokens:  repos.EmailVerificationTokens(),
		UnitOfWork:          unitOfWork,
		PasswordHasher:      passwordHasher,
		AccessTokenProvider: tokenProvider,
		TokenGenerator:      tokenGenerator,
		JobPublisher:        jobPublisher,
		AppBaseURL:          deps.Config.App.BaseURL,
		RefreshTokenTTL:     time.Duration(deps.Config.JWT.RefreshTTLHours) * time.Hour,
	})
	profileService := appService.NewProfileService(repos.Users())
	adminUserService := appService.NewAdminUserService(repos.Users(), repos.Roles(), unitOfWork, passwordHasher)
	roleService := appService.NewRoleService(repos.Roles(), repos.Permissions(), unitOfWork)
	permissionService := appService.NewPermissionService(repos.Permissions())

	return &Module{
		AuthHandler:          identityHandler.NewAuthHandler(authService),
		ProfileHandler:       identityHandler.NewProfileHandler(profileService),
		AdminUserHandler:     identityHandler.NewAdminUserHandler(adminUserService),
		RoleHandler:          identityHandler.NewRoleHandler(roleService),
		PermissionHandler:    identityHandler.NewPermissionHandler(permissionService),
		AuthMiddleware:       identityMiddleware.NewAuthMiddleware(tokenProvider),
		PermissionMiddleware: identityMiddleware.NewPermissionMiddleware(repos.Users()),
	}
}

func (m *Module) RegisterRoutes(router *gin.Engine) {
	identityHTTP.RegisterRoutes(router, identityHTTP.RouterDependencies{
		AuthHandler:          m.AuthHandler,
		ProfileHandler:       m.ProfileHandler,
		AdminUserHandler:     m.AdminUserHandler,
		RoleHandler:          m.RoleHandler,
		PermissionHandler:    m.PermissionHandler,
		AuthMiddleware:       m.AuthMiddleware,
		PermissionMiddleware: m.PermissionMiddleware,
	})
}
