package queue

import (
	"context"
	"time"
)

type Job struct {
	Name        string
	Payload     []byte
	Attempts    int
	MaxAttempts int
	RunAt       time.Time
}

type Handler interface {
	Handle(ctx context.Context, job Job) error
}

type Publisher interface {
	Publish(ctx context.Context, job Job) error
}

type Consumer interface {
	Register(jobName string, handler Handler)
	Start(ctx context.Context) error
}
