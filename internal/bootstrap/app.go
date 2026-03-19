package bootstrap

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/ramdhanrizki/bytecode-api/configs"
	"github.com/ramdhanrizki/bytecode-api/internal/category"
	"github.com/ramdhanrizki/bytecode-api/internal/identity"
	identityService "github.com/ramdhanrizki/bytecode-api/internal/identity/application/service"
	identityWorker "github.com/ramdhanrizki/bytecode-api/internal/identity/delivery/worker"
	platformSMTP "github.com/ramdhanrizki/bytecode-api/internal/platform/mail/smtp"
	"github.com/ramdhanrizki/bytecode-api/internal/product"
	sharedLogger "github.com/ramdhanrizki/bytecode-api/internal/shared/logger"
	sharedQueue "github.com/ramdhanrizki/bytecode-api/internal/shared/queue"
	workerpkg "github.com/ramdhanrizki/bytecode-api/internal/worker"
	workerJobs "github.com/ramdhanrizki/bytecode-api/internal/worker/jobs"
	workerQueue "github.com/ramdhanrizki/bytecode-api/internal/worker/queue"
)

type App struct {
	Config   configs.Config
	Logger   sharedLogger.Logger
	DB       *gorm.DB
	Server   *http.Server
	Queue    sharedQueue.Publisher
	Consumer sharedQueue.Consumer
	Identity *identity.Module
	Category *category.Module
	Product  *product.Module
	Worker   *workerpkg.Server
}

func NewApp() (*App, error) {
	cfg, err := LoadConfig()
	if err != nil {
		return nil, err
	}

	logger, err := NewLogger(cfg)
	if err != nil {
		return nil, fmt.Errorf("initialize logger: %w", err)
	}

	db, err := NewDatabase(cfg, logger)
	if err != nil {
		return nil, fmt.Errorf("initialize database: %w", err)
	}

	queue, err := NewQueue(cfg, logger, db)
	if err != nil {
		return nil, fmt.Errorf("initialize queue publisher: %w", err)
	}

	consumer, err := NewQueueConsumer(cfg, logger, db)
	if err != nil {
		return nil, fmt.Errorf("initialize queue consumer: %w", err)
	}

	identityModule := identity.NewModule(identity.Dependencies{
		Config: cfg,
		Logger: logger,
		DB:     db,
		Queue:  queue,
	})
	categoryModule := category.NewModule(category.Dependencies{Logger: logger, DB: db})
	productModule := product.NewModule(product.Dependencies{Logger: logger, DB: db})

	mailSender := platformSMTP.NewSender(cfg.SMTP, logger)
	emailVerificationService := identityService.NewEmailVerificationDeliveryService(logger, mailSender, cfg.App.Name)
	emailVerificationHandler := identityWorker.NewEmailVerificationHandler(emailVerificationService, logger)
	registry := workerQueue.NewRegistry(consumer, logger)
	workerJobs.Register(registry, workerJobs.Dependencies{
		EmailVerificationHandler: emailVerificationHandler,
	})
	workerServer := workerpkg.NewServer(consumer)

	zapLogger := unwrapZapLogger(logger)
	server := NewHTTPServer(cfg, zapLogger, identityModule, categoryModule, productModule)

	return &App{
		Config:   cfg,
		Logger:   logger,
		DB:       db,
		Server:   server,
		Queue:    queue,
		Consumer: consumer,
		Identity: identityModule,
		Category: categoryModule,
		Product:  productModule,
		Worker:   workerServer,
	}, nil
}

func (a *App) RunAPI() error {
	a.Logger.Info("starting api server", zap.String("port", a.Config.App.Port))

	errCh := make(chan error, 1)
	go func() {
		if err := a.Server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
	}()

	return a.waitForShutdown(errCh, func(ctx context.Context) error {
		return a.Server.Shutdown(ctx)
	})
}

func (a *App) RunWorker() error {
	a.Logger.Info("starting worker", zap.Int("concurrency", a.Config.Worker.Concurrency))

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	errCh := make(chan error, 1)
	go func() {
		if err := a.Worker.Run(ctx); err != nil {
			errCh <- err
		}
	}()

	select {
	case <-ctx.Done():
		a.Logger.Info("worker shutdown signal received")
		return nil
	case err := <-errCh:
		return err
	}
}

func (a *App) waitForShutdown(errCh <-chan error, shutdown func(context.Context) error) error {
	signalCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	select {
	case <-signalCtx.Done():
		a.Logger.Info("shutdown signal received")
	case err := <-errCh:
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), a.Config.App.Shutdown)
	defer cancel()

	if err := shutdown(ctx); err != nil {
		return err
	}

	if a.DB != nil {
		sqlDB, err := a.DB.DB()
		if err == nil {
			_ = sqlDB.Close()
		}
	}

	if err := a.Logger.Sync(); err != nil {
		return nil
	}

	time.Sleep(100 * time.Millisecond)
	return nil
}

func unwrapZapLogger(logger sharedLogger.Logger) *zap.Logger {
	if typed, ok := logger.(*sharedLogger.ZapLogger); ok {
		return typed.Base()
	}

	fallback, _ := zap.NewProduction()
	return fallback
}
