package postgres

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"

	"github.com/ramdhanrizki/bytecode-api/configs"
	sharedLogger "github.com/ramdhanrizki/bytecode-api/internal/shared/logger"
)

func NewGormDB(cfg configs.DatabaseConfig, logger sharedLogger.Logger) (*gorm.DB, error) {
	gormConfig := &gorm.Config{
		Logger: newGormLogger(logger),
	}

	db, err := gorm.Open(postgres.Open(databaseDSN(cfg)), gormConfig)
	if err != nil {
		return nil, fmt.Errorf("open gorm connection: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("access sql db: %w", err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetConnMaxLifetime(30 * time.Minute)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := sqlDB.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("ping database: %w", err)
	}

	logger.Info("database connection initialized", zap.String("host", cfg.Host), zap.Int("port", cfg.Port), zap.String("db", cfg.Name))

	return db, nil
}

func databaseDSN(cfg configs.DatabaseConfig) string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s TimeZone=UTC",
		cfg.Host,
		cfg.Port,
		cfg.User,
		cfg.Password,
		cfg.Name,
		cfg.SSLMode,
	)
}

func newGormLogger(logger sharedLogger.Logger) gormlogger.Interface {
	return gormlogger.New(
		&gormLogWriter{logger: logger},
		gormlogger.Config{
			SlowThreshold:             500 * time.Millisecond,
			LogLevel:                  gormlogger.Warn,
			IgnoreRecordNotFoundError: true,
			Colorful:                  false,
		},
	)
}

type gormLogWriter struct {
	logger sharedLogger.Logger
}

func (w *gormLogWriter) Printf(format string, args ...any) {
	w.logger.Debug(fmt.Sprintf(format, args...))
}
