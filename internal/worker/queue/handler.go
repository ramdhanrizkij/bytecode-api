package queue

import (
	"context"

	sharedQueue "github.com/ramdhanrizki/bytecode-api/internal/shared/queue"
)

type NoopHandler struct{}

func (NoopHandler) Handle(_ context.Context, _ sharedQueue.Job) error {
	return nil
}
