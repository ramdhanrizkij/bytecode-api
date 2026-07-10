package worker

import (
	"context"
	"sync"
	"time"

	"go.uber.org/zap"
)

// ScheduledTask represents a task that runs at a regular interval.
type ScheduledTask struct {
	Name     string
	Interval time.Duration
	Task     func(ctx context.Context) error
}

// Scheduler manages multiple periodic tasks.
type Scheduler struct {
	tasks  []ScheduledTask
	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
	logger *zap.Logger
}

// NewScheduler creates a new Scheduler instance.
func NewScheduler(logger *zap.Logger) *Scheduler {
	ctx, cancel := context.WithCancel(context.Background())
	return &Scheduler{
		tasks:  make([]ScheduledTask, 0),
		ctx:    ctx,
		cancel: cancel,
		logger: logger,
	}
}

// Register adds a task to the scheduler.
func (s *Scheduler) Register(task ScheduledTask) {
	s.tasks = append(s.tasks, task)
}

// Start begins the execution of all registered tasks.
func (s *Scheduler) Start() {
	s.logger.Info("starting scheduler", zap.Int("tasks", len(s.tasks)))
	for _, task := range s.tasks {
		s.wg.Add(1)
		go s.runTask(task)
	}
}

// Stop gracefully stops all scheduled tasks.
func (s *Scheduler) Stop() {
	s.logger.Info("stopping scheduler...")
	s.cancel()
	s.wg.Wait()
	s.logger.Info("scheduler stopped")
}

func (s *Scheduler) runTask(task ScheduledTask) {
	defer s.wg.Done()

	ticker := time.NewTicker(task.Interval)
	defer ticker.Stop()

	s.logger.Debug("task scheduled", 
		zap.String("task_name", task.Name), 
		zap.Duration("interval", task.Interval))

	for {
		select {
		case <-ticker.C:
			s.executeTask(task)
		case <-s.ctx.Done():
			return
		}
	}
}

func (s *Scheduler) executeTask(task ScheduledTask) {
	defer func() {
		if r := recover(); r != nil {
			s.logger.Error("scheduled task panic", 
				zap.String("task_name", task.Name), 
				zap.Any("panic", r))
		}
	}()

	start := time.Now()
	err := task.Task(s.ctx)
	duration := time.Since(start)

	if err != nil {
		s.logger.Error("scheduled task error",
			zap.String("task_name", task.Name),
			zap.Error(err),
			zap.Duration("duration", duration))
	} else {
		s.logger.Debug("scheduled task completed",
			zap.String("task_name", task.Name),
			zap.Duration("duration", duration))
	}
}
