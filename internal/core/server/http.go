package server

import (
	"errors"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/ramdhanrizkij/bytecode-api/internal/core/config"
	"github.com/ramdhanrizkij/bytecode-api/internal/core/middleware"
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
	config    *config.Config
	logger    *zap.Logger
	worker    *worker.WorkerPool
	scheduler *worker.Scheduler
}

// NewServer initializes the Fiber application and its global settings.
func NewServer(cfg *config.Config, db *gorm.DB, logger *zap.Logger, wp *worker.WorkerPool, sched *worker.Scheduler) *Server {
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

	return &Server{
		app:       app,
		db:        db,
		config:    cfg,
		logger:    logger,
		worker:    wp,
		scheduler: sched,
	}
}

// SetupRoutes handles the Dependency Injection and route registration for all modules.
func (s *Server) SetupRoutes() {
	api := s.app.Group("/api/v1")

	// 1. Feature: Auth
	authRepository := authRepo.NewAuthRepository(s.db)
	authServ := authService.NewAuthService(authRepository, s.worker, s.config.JWT.Secret, s.config.JWT.ExpiryHours, s.logger)
	authHdl := authHandler.NewAuthHTTPHandler(authServ, s.logger)
	authHandler.RegisterRoutes(api, authHdl)

	// 2. Feature: Role
	roleRepository := roleRepo.NewRoleRepository(s.db)
	roleServ := roleService.NewRoleService(roleRepository, s.logger)
	roleHdl := roleHandler.NewRoleHTTPHandler(roleServ, s.logger)
	roleHandler.RegisterRoutes(api, roleHdl, s.db, s.config.JWT.Secret)

	// 3. Feature: Permission
	permRepository := permRepo.NewPermissionRepository(s.db)
	permServ := permService.NewPermissionService(permRepository, s.logger)
	permHdl := permHandler.NewPermissionHTTPHandler(permServ, s.logger)
	permHandler.RegisterRoutes(api, permHdl, s.db, s.config.JWT.Secret)

	// 4. Feature: User
	userRepository := userRepo.NewUserRepository(s.db)
	userServ := userService.NewUserService(userRepository, s.logger)
	userHdl := userHandler.NewUserHTTPHandler(userServ, s.logger)
	userHandler.RegisterRoutes(api, userHdl, s.db, s.config.JWT.Secret)

	// Catch-all route for 404 Not Found
	s.app.Use(func(c *fiber.Ctx) error {
		return response.Error(c, fiber.StatusNotFound, "route not found")
	})
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
	return s.app.Shutdown()
}

// AppForTest returns the underlying Fiber app instance for testing purposes.
func (s *Server) AppForTest() *fiber.App {
	return s.app
}

// customErrorHandler converts AppError into the standard API response format.
func customErrorHandler(c *fiber.Ctx, err error) error {
	var appErr *apperrors.AppError
	if errors.As(err, &appErr) {
		return response.Error(c, appErr.Code, appErr.Message)
	}

	return response.Error(c, fiber.StatusInternalServerError, "internal server error")
}
