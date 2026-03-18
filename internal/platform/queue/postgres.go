package queue

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"

	sharedLogger "github.com/ramdhanrizki/bytecode-api/internal/shared/logger"
	sharedQueue "github.com/ramdhanrizki/bytecode-api/internal/shared/queue"
)

const (
	defaultMaxAttempts  = 5
	defaultPollInterval = 1 * time.Second
	reservationTimeout  = 5 * time.Minute
	maxLastErrorLength  = 2000
)

var errNoAvailableJob = errors.New("no available jobs")

type PostgresQueue struct {
	db           *gorm.DB
	logger       sharedLogger.Logger
	concurrency  int
	pollInterval time.Duration
	handlers     map[string]sharedQueue.Handler
	mu           sync.RWMutex
}

type workerJobModel struct {
	ID          string     `gorm:"column:id;type:uuid;primaryKey"`
	Name        string     `gorm:"column:name;index"`
	Payload     []byte     `gorm:"column:payload;type:bytea"`
	Attempts    int        `gorm:"column:attempts;not null;default:0"`
	MaxAttempts int        `gorm:"column:max_attempts;not null;default:5"`
	RunAt       time.Time  `gorm:"column:run_at;index"`
	ReservedAt  *time.Time `gorm:"column:reserved_at;index"`
	CompletedAt *time.Time `gorm:"column:completed_at;index"`
	FailedAt    *time.Time `gorm:"column:failed_at;index"`
	LastError   *string    `gorm:"column:last_error;type:text"`
	CreatedAt   time.Time  `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt   time.Time  `gorm:"column:updated_at;autoUpdateTime"`
}

func (workerJobModel) TableName() string {
	return "worker_jobs"
}

func NewPostgresQueue(db *gorm.DB, logger sharedLogger.Logger, concurrency int) (*PostgresQueue, error) {
	if concurrency <= 0 {
		concurrency = 1
	}

	return &PostgresQueue{
		db:           db,
		logger:       logger,
		concurrency:  concurrency,
		pollInterval: defaultPollInterval,
		handlers:     make(map[string]sharedQueue.Handler),
	}, nil
}

func (q *PostgresQueue) Publish(ctx context.Context, job sharedQueue.Job) error {
	now := time.Now().UTC()
	runAt := job.RunAt
	if runAt.IsZero() {
		runAt = now
	}
	maxAttempts := job.MaxAttempts
	if maxAttempts <= 0 {
		maxAttempts = defaultMaxAttempts
	}

	model := workerJobModel{
		ID:          uuid.NewString(),
		Name:        job.Name,
		Payload:     job.Payload,
		Attempts:    job.Attempts,
		MaxAttempts: maxAttempts,
		RunAt:       runAt,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := q.db.WithContext(ctx).Create(&model).Error; err != nil {
		return err
	}

	q.logger.Info("job published",
		zap.String("job_name", job.Name),
		zap.String("job_id", model.ID),
		zap.Int("max_attempts", maxAttempts),
	)

	return nil
}

func (q *PostgresQueue) Register(jobName string, handler sharedQueue.Handler) {
	q.mu.Lock()
	defer q.mu.Unlock()

	q.handlers[jobName] = handler
}

func (q *PostgresQueue) Start(ctx context.Context) error {
	q.logger.Info("queue consumer started", zap.Int("concurrency", q.concurrency))

	var wg sync.WaitGroup
	for workerIndex := 0; workerIndex < q.concurrency; workerIndex++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			q.runWorker(ctx, index+1)
		}(workerIndex)
	}

	<-ctx.Done()
	wg.Wait()
	q.logger.Info("queue consumer stopped")
	return nil
}

func (q *PostgresQueue) runWorker(ctx context.Context, workerNumber int) {
	for {
		if ctx.Err() != nil {
			return
		}

		job, err := q.claimNextJob(ctx)
		if err != nil {
			if errors.Is(err, errNoAvailableJob) {
				if !sleepWithContext(ctx, q.pollInterval) {
					return
				}
				continue
			}

			q.logger.Error("failed to claim job", zap.Int("worker", workerNumber), zap.Error(err))
			if !sleepWithContext(ctx, q.pollInterval) {
				return
			}
			continue
		}

		if err := q.handleClaimedJob(ctx, workerNumber, job); err != nil {
			q.logger.Error("failed to finalize job handling", zap.Int("worker", workerNumber), zap.String("job_name", job.Name), zap.String("job_id", job.ID), zap.Error(err))
		}
	}
}

func (q *PostgresQueue) handleClaimedJob(ctx context.Context, workerNumber int, job workerJobModel) error {
	q.logger.Info("processing job",
		zap.Int("worker", workerNumber),
		zap.String("job_name", job.Name),
		zap.String("job_id", job.ID),
		zap.Int("attempt", job.Attempts),
	)

	err := q.dispatch(ctx, job)
	if err == nil {
		if completeErr := q.markCompleted(ctx, job); completeErr != nil {
			return completeErr
		}

		q.logger.Info("job completed",
			zap.Int("worker", workerNumber),
			zap.String("job_name", job.Name),
			zap.String("job_id", job.ID),
		)
		return nil
	}

	if retryErr := q.markFailedAttempt(ctx, job, err); retryErr != nil {
		return retryErr
	}

	q.logger.Error("job execution failed",
		zap.Int("worker", workerNumber),
		zap.String("job_name", job.Name),
		zap.String("job_id", job.ID),
		zap.Int("attempt", job.Attempts),
		zap.Error(err),
	)

	return nil
}

func (q *PostgresQueue) dispatch(ctx context.Context, job workerJobModel) error {
	q.mu.RLock()
	handler, ok := q.handlers[job.Name]
	q.mu.RUnlock()
	if !ok {
		return errors.New("queue handler not registered")
	}

	return handler.Handle(ctx, sharedQueue.Job{
		Name:        job.Name,
		Payload:     job.Payload,
		Attempts:    job.Attempts,
		MaxAttempts: job.MaxAttempts,
		RunAt:       job.RunAt,
	})
}

func (q *PostgresQueue) claimNextJob(ctx context.Context) (workerJobModel, error) {
	now := time.Now().UTC()
	staleReservationBefore := now.Add(-reservationTimeout)

	var claimed workerJobModel
	err := q.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		result := tx.Raw(`
			SELECT id, name, payload, attempts, max_attempts, run_at, reserved_at, completed_at, failed_at, last_error, created_at, updated_at
			FROM worker_jobs
			WHERE completed_at IS NULL
			  AND failed_at IS NULL
			  AND attempts < max_attempts
			  AND run_at <= ?
			  AND (reserved_at IS NULL OR reserved_at <= ?)
			ORDER BY run_at ASC, created_at ASC
			LIMIT 1
			FOR UPDATE SKIP LOCKED
		`, now, staleReservationBefore).Scan(&claimed)
		if result.Error != nil {
			return result.Error
		}
		if claimed.ID == "" {
			return errNoAvailableJob
		}

		claimed.Attempts++
		claimed.ReservedAt = &now
		claimed.UpdatedAt = now

		return tx.Model(&workerJobModel{}).
			Where("id = ?", claimed.ID).
			Updates(map[string]any{
				"attempts":    claimed.Attempts,
				"reserved_at": now,
				"updated_at":  now,
			}).Error
	})
	if err != nil {
		return workerJobModel{}, err
	}

	return claimed, nil
}

func (q *PostgresQueue) markCompleted(ctx context.Context, job workerJobModel) error {
	now := time.Now().UTC()
	return q.db.WithContext(ctx).Model(&workerJobModel{}).
		Where("id = ?", job.ID).
		Updates(map[string]any{
			"completed_at": now,
			"reserved_at":  nil,
			"last_error":   nil,
			"updated_at":   now,
		}).Error
}

func (q *PostgresQueue) markFailedAttempt(ctx context.Context, job workerJobModel, cause error) error {
	now := time.Now().UTC()
	message := truncateError(cause)
	updates := map[string]any{
		"reserved_at": nil,
		"last_error":  message,
		"updated_at":  now,
	}

	if job.Attempts >= job.MaxAttempts {
		updates["failed_at"] = now
		return q.db.WithContext(ctx).Model(&workerJobModel{}).Where("id = ?", job.ID).Updates(updates).Error
	}

	updates["run_at"] = now.Add(backoffForAttempt(job.Attempts))
	return q.db.WithContext(ctx).Model(&workerJobModel{}).Where("id = ?", job.ID).Updates(updates).Error
}

func backoffForAttempt(attempt int) time.Duration {
	if attempt <= 1 {
		return 5 * time.Second
	}
	if attempt == 2 {
		return 15 * time.Second
	}
	if attempt == 3 {
		return 30 * time.Second
	}
	return time.Minute
}

func sleepWithContext(ctx context.Context, delay time.Duration) bool {
	timer := time.NewTimer(delay)
	defer timer.Stop()

	select {
	case <-ctx.Done():
		return false
	case <-timer.C:
		return true
	}
}

func truncateError(err error) string {
	if err == nil {
		return ""
	}
	message := err.Error()
	if len(message) > maxLastErrorLength {
		return message[:maxLastErrorLength]
	}
	return message
}
