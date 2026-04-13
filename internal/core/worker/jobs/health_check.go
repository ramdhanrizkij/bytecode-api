package jobs

import (
	"context"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

// HealthCheckJob periodically checks the system health, primarily database connectivity.
type HealthCheckJob struct {
	db  *gorm.DB
	log *zap.Logger
}

// NewHealthCheckJob creates a new HealthCheckJob.
func NewHealthCheckJob(db *gorm.DB, log *zap.Logger) *HealthCheckJob {
	return &HealthCheckJob{
		db:  db,
		log: log,
	}
}

// Name returns the job name.
func (j *HealthCheckJob) Name() string {
	return "system_health_check"
}

// Execute checks the database connection.
func (j *HealthCheckJob) Execute(ctx context.Context) error {
	sqlDB, err := j.db.DB()
	if err != nil {
		return err
	}

	if err := sqlDB.PingContext(ctx); err != nil {
		j.log.Error("system health check: database unreachable", zap.Error(err))
		return err
	}

	j.log.Info("system health: OK")
	return nil
}
