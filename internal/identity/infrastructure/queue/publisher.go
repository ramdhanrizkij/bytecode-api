package queue

import (
	"context"
	"encoding/json"

	"github.com/ramdhanrizki/bytecode-api/internal/identity/application/job"
	sharedQueue "github.com/ramdhanrizki/bytecode-api/internal/shared/queue"
)

type VerificationPublisher struct {
	queue sharedQueue.Publisher
}

func NewVerificationPublisher(queue sharedQueue.Publisher) *VerificationPublisher {
	return &VerificationPublisher{queue: queue}
}

func (p *VerificationPublisher) PublishEmailVerification(ctx context.Context, payload job.EmailVerificationJob) error {
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	return p.queue.Publish(ctx, sharedQueue.Job{
		Name:        job.EmailVerificationJobName,
		Payload:     body,
		Attempts:    0,
		MaxAttempts: 5,
	})
}
