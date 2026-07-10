package server

import (
	"context"
	"errors"
	"fmt"
	"time"

	fiberSwagger "github.com/gofiber/contrib/v3/swaggo"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/basicauth"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/recover"
	"github.com/gofiber/fiber/v3/middleware/static"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/ramdhanrizkij/bytecode-api/internal/core/cache"
	"github.com/ramdhanrizkij/bytecode-api/internal/core/config"
	"github.com/ramdhanrizkij/bytecode-api/internal/core/middleware"
	"github.com/ramdhanrizkij/bytecode-api/internal/core/storage"
	"github.com/ramdhanrizkij/bytecode-api/internal/core/worker"
	apperrors "github.com/ramdhanrizkij/bytecode-api/internal/shared/errors"
	"github.com/ramdhanrizkij/bytecode-api/internal/shared/response"

	// Feature Auth
	authHandler "github.com/ramdhanrizkij/bytecode-api/internal/features/auth/handler"
	authRepo "github.com/ramdhanrizkij/bytecode-api/internal/features/auth/repository"
	authService "github.com/ramdhanrizkij/bytecode-api/internal/features/auth/service"

	// Feature Role
	roleHandler "github.com/ramdhanrizkij/bytecode-api/internal/features/role/handler"
	roleRepo "github.com/ramdhanrizkij/bytecode-api/internal/features/role/repository"
	roleService "github.com/ramdhanrizkij/bytecode-api/internal/features/role/service"

	// Feature Permission
	permHandler "github.com/ramdhanrizkij/bytecode-api/internal/features/permission/handler"
	permRepo "github.com/ramdhanrizkij/bytecode-api/internal/features/permission/repository"
	permService "github.com/ramdhanrizkij/bytecode-api/internal/features/permission/service"

	// Feature User
	userHandler "github.com/ramdhanrizkij/bytecode-api/internal/features/user/handler"
	userRepo "github.com/ramdhanrizkij/bytecode-api/internal/features/user/repository"
	userService "github.com/ramdhanrizkij/bytecode-api/internal/features/user/service"
)

// Server represents the HTTP server container.
type Server struct {
	app       *fiber.App
	db        *gorm.DB
	cache     cache.Client
	storage   storage.Provider
	config    *config.Config
	logger    *zap.Logger
	worker    *worker.WorkerPool
	scheduler *worker.Scheduler
}

// NewServer initializes the Fiber application and its global settings.
func NewServer(cfg *config.Config, db *gorm.DB, logger *zap.Logger, wp *worker.WorkerPool, sched *worker.Scheduler) (*Server, error) {
	app := fiber.New(fiber.Config{
		AppName:      cfg.App.Name,
		BodyLimit:    4 * 1024 * 1024, // 4MB
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		ErrorHandler: customErrorHandler,
	})

	// Global Middleware
	app.Use(recover.New())
	app.Use(cors.New())
	app.Use(middleware.RequestLogger(logger))

	cacheClient, err := cache.NewClient(&cfg.Redis, logger)
	if err != nil {
		return nil, err
	}

	storageProvider, err := storage.NewProvider(&cfg.Storage, logger)
	if err != nil {
		return nil, err
	}

	return &Server{
		app:       app,
		db:        db,
		cache:     cacheClient,
		storage:   storageProvider,
		config:    cfg,
		logger:    logger,
		worker:    wp,
		scheduler: sched,
	}, nil
}

// SetupRoutes handles the Dependency Injection and route registration for all modules.
func (s *Server) SetupRoutes() {
	if s.storage.ProviderName() == storage.ProviderLocal {
		s.app.Get(s.storageBaseURL()+"/*", static.New(s.storageLocalPath()))
	}

	s.registerSwagger()

	api := s.app.Group("/api/v1")
	api.Get("/health", s.healthCheck)

	// 1. Feature: Auth
	authRepository := authRepo.NewAuthRepository(s.db)
	authServ := authService.NewAuthService(
		authRepository,
		s.worker,
		s.config.JWT.Secret,
		s.config.JWT.ExpiryHours,
		s.config.JWT.RefreshExpiryHours,
		s.logger,
	)
	authHdl := authHandler.NewAuthHTTPHandler(authServ, s.logger)
	authHandler.RegisterRoutes(api, authHdl)

	// 2. Feature: Role
	roleRepository := roleRepo.NewRoleRepository(s.db)
	roleServ := roleService.NewRoleService(roleRepository, s.cache, s.cacheTTL(), s.logger)
	roleHdl := roleHandler.NewRoleHTTPHandler(roleServ, s.logger)
	roleHandler.RegisterRoutes(api, roleHdl, s.db, s.config.JWT.Secret)

	// 3. Feature: Permission
	permRepository := permRepo.NewPermissionRepository(s.db)
	permServ := permService.NewPermissionService(permRepository, s.cache, s.cacheTTL(), s.logger)
	permHdl := permHandler.NewPermissionHTTPHandler(permServ, s.logger)
	permHandler.RegisterRoutes(api, permHdl, s.db, s.config.JWT.Secret)

	// 4. Feature: User
	userRepository := userRepo.NewUserRepository(s.db)
	userServ := userService.NewUserService(
		userRepository,
		s.cache,
		s.storage,
		s.config.Storage.DefaultBucket,
		s.cacheTTL(),
		s.logger,
	)
	userHdl := userHandler.NewUserHTTPHandler(userServ, s.logger)
	userHandler.RegisterRoutes(api, userHdl, s.db, s.config.JWT.Secret)

	// Catch-all route for 404 Not Found
	s.app.Use(func(c fiber.Ctx) error {
		return response.Error(c, fiber.StatusNotFound, "route not found")
	})
}

func (s *Server) registerSwagger() {
	if !s.config.Swagger.Enabled {
		s.logger.Info("swagger documentation disabled")
		return
	}

	if s.config.App.IsProduction() {
		if s.config.Swagger.Username == "" || s.config.Swagger.Password == "" {
			s.logger.Warn("swagger documentation disabled in production because basic auth credentials are missing")
			return
		}

		swaggerGroup := s.app.Group("/swagger", basicauth.New(basicauth.Config{
			Users: map[string]string{
				s.config.Swagger.Username: s.config.Swagger.Password,
			},
		}))
		swaggerGroup.Get("/*", fiberSwagger.New())
		return
	}

	s.app.Get("/swagger/*", fiberSwagger.New())
}

// Start runs the HTTP server.
func (s *Server) Start() error {
	addr := fmt.Sprintf(":%d", s.config.App.Port)
	s.logger.Info("server starting", zap.Int("port", s.config.App.Port))
	return s.app.Listen(addr)
}

// Shutdown gracefully stops the HTTP server.
func (s *Server) Shutdown() error {
	s.logger.Info("shutting down HTTP server...")
	if err := s.app.Shutdown(); err != nil {
		return err
	}
	if err := s.storage.Close(); err != nil {
		return err
	}
	return s.cache.Close()
}

// AppForTest returns the underlying Fiber app instance for testing purposes.
func (s *Server) AppForTest() *fiber.App {
	return s.app
}

// healthCheck godoc
// @Summary Check service health
// @Description Returns the current health state for the API, database, cache, and storage provider.
// @Tags Health
// @Produce json
// @Success 200 {object} swaggerdocs.HealthResponse
// @Router /health [get]
func (s *Server) healthCheck(c fiber.Ctx) error {
	dbStatus := "up"
	cacheStatus := "disabled"

	if s.cache.IsEnabled() {
		cacheStatus = "up"
	}

	sqlDB, err := s.db.DB()
	if err != nil {
		s.logger.Error("failed to get sql.DB for health check", zap.Error(err))
		dbStatus = "down"
	} else {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		if err := sqlDB.PingContext(ctx); err != nil {
			s.logger.Error("database ping failed during health check", zap.Error(err))
			dbStatus = "down"
		}
	}

	status := "ok"
	message := "service is healthy"
	if dbStatus != "up" {
		status = "degraded"
		message = "service is unhealthy"
	}

	return c.Status(fiber.StatusOK).JSON(response.Response{
		Meta: response.Meta{
			Code:    fiber.StatusOK,
			Message: message,
		},
		Data: fiber.Map{
			"status":      status,
			"service":     s.config.App.Name,
			"environment": s.config.App.Env,
			"database":    dbStatus,
			"cache":       cacheStatus,
			"storage":     s.storage.ProviderName(),
		},
	})
}

func (s *Server) cacheTTL() time.Duration {
	ttlMinutes := s.config.Redis.CacheTTLMinutes
	if ttlMinutes <= 0 {
		ttlMinutes = 5
	}
	return time.Duration(ttlMinutes) * time.Minute
}

func (s *Server) storageBaseURL() string {
	if s.config.Storage.BaseURL == "" {
		return "/storage"
	}
	return s.config.Storage.BaseURL
}

func (s *Server) storageLocalPath() string {
	if s.config.Storage.LocalPath == "" {
		return "storage"
	}
	return s.config.Storage.LocalPath
}

// customErrorHandler converts AppError into the standard API response format.
func customErrorHandler(c fiber.Ctx, err error) error {
	var appErr *apperrors.AppError
	if errors.As(err, &appErr) {
		return response.Error(c, appErr.Code, appErr.Message)
	}

	return response.Error(c, fiber.StatusInternalServerError, "internal server error")
}
