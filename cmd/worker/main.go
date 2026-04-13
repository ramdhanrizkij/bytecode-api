package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"

	"github.com/ramdhanrizkij/bytecode-api/internal/core/config"
	"github.com/ramdhanrizkij/bytecode-api/internal/core/database"
	"github.com/ramdhanrizkij/bytecode-api/internal/core/worker"
	"github.com/ramdhanrizkij/bytecode-api/internal/core/worker/jobs"
	authRepo "github.com/ramdhanrizkij/bytecode-api/internal/features/auth/repository"
	authService "github.com/ramdhanrizkij/bytecode-api/internal/features/auth/service"
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

	// 4. Initialize Scheduler
	sched := worker.NewScheduler(logger.Log)

	// Job 1: Health Check (every 5 minutes)
	healthJob := jobs.NewHealthCheckJob(db, logger.Log)
	sched.Register(worker.ScheduledTask{
		Name:     healthJob.Name(),
		Interval: 5 * time.Minute,
		Task:     healthJob.Execute,
	})

	// Job 2: Auth Token Cleanup (every 1 hour)
	// We need a dummy worker pool here because NewAuthService requires it, 
	// even though it's not used for CleanupExpiredTokens.
	dummyWP := worker.NewWorkerPool(1, 1, logger.Log) 
	authRepository := authRepo.NewAuthRepository(db)
	authServ := authService.NewAuthService(authRepository, dummyWP, cfg.JWT.Secret, cfg.JWT.ExpiryHours, logger.Log)
	
	sched.Register(worker.ScheduledTask{
		Name:     "cleanup_expired_tokens",
		Interval: 1 * time.Hour,
		Task:     authServ.CleanupExpiredTokens,
	})

	// 5. Start Scheduler
	sched.Start()
	logger.Log.Info("worker process started successfully")

	// 6. Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Log.Info("shutting down worker process...")
	sched.Stop()

	if err := database.CloseDB(db); err != nil {
		logger.Log.Error("failed to close database connection", zap.Error(err))
	}

	logger.Log.Info("worker process stopped gracefully")
}
