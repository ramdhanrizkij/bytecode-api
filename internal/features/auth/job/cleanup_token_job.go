package job

import (
	"context"

	"github.com/ramdhanrizkij/bytecode-api/internal/features/auth/domain"
)

// CleanupExpiredTokensJob implements worker.Job to periodicly clean up auth tokens.
type CleanupExpiredTokensJob struct {
	authService domain.AuthService
}

// NewCleanupExpiredTokensJob creates a new CleanupExpiredTokensJob.
func NewCleanupExpiredTokensJob(authService domain.AuthService) *CleanupExpiredTokensJob {
	return &CleanupExpiredTokensJob{
		authService: authService,
	}
}

// Name returns the job name.
func (j *CleanupExpiredTokensJob) Name() string {
	return "cleanup_expired_tokens"
}

// Execute runs the cleanup logic via the auth service.
func (j *CleanupExpiredTokensJob) Execute(ctx context.Context) error {
	return j.authService.CleanupExpiredTokens(ctx)
}
