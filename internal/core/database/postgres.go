package database

import (
	"fmt"
	"time"

	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"

	"github.com/ramdhanrizkij/bytecode-api/internal/core/config"
)

// NewPostgresDB initializes a GORM database connection to PostgreSQL.
// It configures the connection pool and GORM logger based on the application environment.
// AutoMigrate is intentionally NOT used — migrations are managed by golang-migrate.
func NewPostgresDB(cfg *config.DBConfig, appEnv string, log *zap.Logger) (*gorm.DB, error) {
	// Choose GORM log level: Info for development, Silent for production.
	gormLogLevel := gormlogger.Info
	if appEnv == "production" {
		gormLogLevel = gormlogger.Silent
	}

	gormCfg := &gorm.Config{
		Logger: gormlogger.Default.LogMode(gormLogLevel),
	}

	db, err := gorm.Open(postgres.Open(cfg.DSN()), gormCfg)
	if err != nil {
		log.Error("failed to connect to database",
			zap.String("host", cfg.Host),
			zap.Int("port", cfg.Port),
			zap.String("dbname", cfg.Name),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	// Configure the underlying *sql.DB connection pool.
	sqlDB, err := db.DB()
	if err != nil {
		log.Error("failed to get underlying sql.DB from GORM", zap.Error(err))
		return nil, fmt.Errorf("failed to get sql.DB: %w", err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(1 * time.Hour)

	// Verify connectivity with a Ping.
	if err := sqlDB.Ping(); err != nil {
		log.Error("database ping failed",
			zap.String("host", cfg.Host),
			zap.Int("port", cfg.Port),
			zap.Error(err),
		)
		return nil, fmt.Errorf("database ping failed: %w", err)
	}

	log.Info("database connection established",
		zap.String("host", cfg.Host),
		zap.Int("port", cfg.Port),
		zap.String("dbname", cfg.Name),
	)

	return db, nil
}

// CloseDB gracefully closes the underlying database connection pool.
func CloseDB(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get sql.DB for closing: %w", err)
	}
	if err := sqlDB.Close(); err != nil {
		return fmt.Errorf("failed to close database connection: %w", err)
	}
	return nil
}
