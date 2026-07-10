package worker

import (
	"context"
	"fmt"
	"sync"
	"time"

	"go.uber.org/zap"
)

// Job defines the interface for a task that can be executed by the worker pool.
type Job interface {
	Name() string
	Execute(ctx context.Context) error
}

// WorkerPool manages a pool of goroutines to execute jobs concurrently.
type WorkerPool struct {
	workers    int
	jobQueue   chan Job
	wg         sync.WaitGroup
	ctx        context.Context
	cancel     context.CancelFunc
	logger     *zap.Logger
}

// NewWorkerPool creates a new WorkerPool with the specified number of workers and queue size.
func NewWorkerPool(workers int, queueSize int, logger *zap.Logger) *WorkerPool {
	ctx, cancel := context.WithCancel(context.Background())
	return &WorkerPool{
		workers:  workers,
		jobQueue: make(chan Job, queueSize),
		ctx:      ctx,
		cancel:   cancel,
		logger:   logger,
	}
}

// Start initiates the worker goroutines.
func (wp *WorkerPool) Start() {
	wp.logger.Info("starting worker pool", zap.Int("workers", wp.workers))
	for i := 0; i < wp.workers; i++ {
		wp.wg.Add(1)
		go wp.worker(i)
	}
}

// Submit adds a job to the queue. Returns an error if the queue is full or context is cancelled.
func (wp *WorkerPool) Submit(job Job) error {
	select {
	case wp.jobQueue <- job:
		return nil
	case <-wp.ctx.Done():
		return wp.ctx.Err()
	default:
		return fmt.Errorf("worker pool queue is full")
	}
}

// Stop gracefully shuts down the worker pool, waiting for all active jobs to finish.
func (wp *WorkerPool) Stop() {
	wp.logger.Info("stopping worker pool...")
	wp.cancel()
	close(wp.jobQueue)
	wp.wg.Wait()
	wp.logger.Info("worker pool stopped")
}

func (wp *WorkerPool) worker(id int) {
	defer wp.wg.Done()
	
	for job := range wp.jobQueue {
		wp.executeJob(id, job)
	}
}

func (wp *WorkerPool) executeJob(workerID int, job Job) {
	defer func() {
		if r := recover(); r != nil {
			wp.logger.Error("worker recovered from panic",
				zap.Int("worker_id", workerID),
				zap.Any("panic", r),
				zap.String("job_name", job.Name()))
		}
	}()

	start := time.Now()
	wp.logger.Debug("starting job", 
		zap.Int("worker_id", workerID), 
		zap.String("job_name", job.Name()))

	err := job.Execute(wp.ctx)
	
	duration := time.Since(start)
	if err != nil {
		wp.logger.Error("job failed",
			zap.Int("worker_id", workerID),
			zap.String("job_name", job.Name()),
			zap.Error(err),
			zap.Duration("duration", duration))
	} else {
		wp.logger.Info("job completed",
			zap.Int("worker_id", workerID),
			zap.String("job_name", job.Name()),
			zap.Duration("duration", duration))
	}
}
