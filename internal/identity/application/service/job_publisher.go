package service

import (
	"context"

	"github.com/ramdhanrizki/bytecode-api/internal/identity/application/job"
)

type VerificationJobPublisher interface {
	PublishEmailVerification(ctx context.Context, payload job.EmailVerificationJob) error
}
