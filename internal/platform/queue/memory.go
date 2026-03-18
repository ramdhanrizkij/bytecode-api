package queue

import (
	"context"
	"errors"
	"sync"

	"go.uber.org/zap"

	sharedLogger "github.com/ramdhanrizki/bytecode-api/internal/shared/logger"
	sharedQueue "github.com/ramdhanrizki/bytecode-api/internal/shared/queue"
)

type InMemoryQueue struct {
	logger   sharedLogger.Logger
	buffer   chan sharedQueue.Job
	handlers map[string]sharedQueue.Handler
	mu       sync.RWMutex
}

func NewInMemoryQueue(logger sharedLogger.Logger, bufferSize int) *InMemoryQueue {
	if bufferSize <= 0 {
		bufferSize = 100
	}

	return &InMemoryQueue{
		logger:   logger,
		buffer:   make(chan sharedQueue.Job, bufferSize),
		handlers: make(map[string]sharedQueue.Handler),
	}
}

func (q *InMemoryQueue) Publish(ctx context.Context, job sharedQueue.Job) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case q.buffer <- job:
		q.logger.Info("job published", zap.String("job_name", job.Name), zap.Int("attempts", job.Attempts))
		return nil
	}
}

func (q *InMemoryQueue) Register(jobName string, handler sharedQueue.Handler) {
	q.mu.Lock()
	defer q.mu.Unlock()

	q.handlers[jobName] = handler
}

func (q *InMemoryQueue) Start(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			q.logger.Info("queue consumer stopped")
			return nil
		case job := <-q.buffer:
			if err := q.dispatch(ctx, job); err != nil {
				q.logger.Error("job dispatch failed", zap.String("job_name", job.Name), zap.Error(err))
			}
		}
	}
}

func (q *InMemoryQueue) dispatch(ctx context.Context, job sharedQueue.Job) error {
	q.mu.RLock()
	handler, ok := q.handlers[job.Name]
	q.mu.RUnlock()
	if !ok {
		return errors.New("queue handler not registered")
	}

	return handler.Handle(ctx, job)
}
