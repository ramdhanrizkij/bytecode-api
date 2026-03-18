package queue

import (
	"go.uber.org/zap"

	sharedLogger "github.com/ramdhanrizki/bytecode-api/internal/shared/logger"
	sharedQueue "github.com/ramdhanrizki/bytecode-api/internal/shared/queue"
)

type Registry struct {
	consumer sharedQueue.Consumer
	logger   sharedLogger.Logger
}

func NewRegistry(consumer sharedQueue.Consumer, logger sharedLogger.Logger) *Registry {
	return &Registry{consumer: consumer, logger: logger}
}

func (r *Registry) Register(jobName string, handler sharedQueue.Handler) {
	r.consumer.Register(jobName, handler)
	r.logger.Info("worker handler registered", zap.String("job_name", jobName))
}
