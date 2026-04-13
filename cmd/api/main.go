package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"

	"github.com/ramdhanrizkij/bytecode-api/internal/core/config"
	"github.com/ramdhanrizkij/bytecode-api/internal/core/database"
	"github.com/ramdhanrizkij/bytecode-api/internal/core/server"
	"github.com/ramdhanrizkij/bytecode-api/internal/core/worker"
	"github.com/ramdhanrizkij/bytecode-api/pkg/logger"
)

func main() {
	// 1. Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// 2. Initialize logger
	if err := logger.InitGlobal(cfg.Log.Level); err != nil {
		log.Fatalf("failed to initialize logger: %v", err)
	}
	defer logger.Log.Sync()

	// 3. Connect to database
	db, err := database.NewPostgresDB(&cfg.DB, cfg.App.Env, logger.Log)
	if err != nil {
		logger.Log.Fatal("failed to connect to database", zap.Error(err))
	}

	// 4. Initialize Worker Pool (for immediate async tasks like sending emails)
	wp := worker.NewWorkerPool(5, 100, logger.Log)
	wp.Start()

	// 5. Initialize Server (Scheduler is passed as nil as it's now a separate process)
	srv := server.NewServer(cfg, db, logger.Log, wp, nil)
	srv.SetupRoutes()

	// 6. Start server in a goroutine
	go func() {
		if err := srv.Start(); err != nil {
			logger.Log.Fatal("failed to start server", zap.Error(err))
		}
	}()

	// 7. Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Log.Info("shutting down API server...")

	// Sequence: Stop Worker Pool -> Shutdown Server -> Close DB
	wp.Stop()
	if err := srv.Shutdown(); err != nil {
		logger.Log.Error("failed to shutdown server gracefully", zap.Error(err))
	}

	// Close database connection
	if err := database.CloseDB(db); err != nil {
		logger.Log.Error("failed to close database connection", zap.Error(err))
	}

	logger.Log.Info("API server stopped gracefully")
}
