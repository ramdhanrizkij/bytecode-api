package job

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"
)

// SendWelcomeEmailJob implements worker.Job to simulate sending a welcome email.
type SendWelcomeEmailJob struct {
	UserEmail string
	UserName  string
	Log       *zap.Logger
}

// NewSendWelcomeEmailJob creates a new instance of SendWelcomeEmailJob.
func NewSendWelcomeEmailJob(email, name string, log *zap.Logger) *SendWelcomeEmailJob {
	return &SendWelcomeEmailJob{
		UserEmail: email,
		UserName:  name,
		Log:       log,
	}
}

// Name returns the unique identifier for this job type.
func (j *SendWelcomeEmailJob) Name() string {
	return fmt.Sprintf("send_welcome_email:%s", j.UserEmail)
}

// Execute simulates the email sending process.
func (j *SendWelcomeEmailJob) Execute(ctx context.Context) error {
	j.Log.Info("attempting to send welcome email", 
		zap.String("email", j.UserEmail), 
		zap.String("name", j.UserName))

	// Simulate network latency or SMTP server interaction
	select {
	case <-time.After(2 * time.Second):
		j.Log.Info("welcome email sent successfully", zap.String("email", j.UserEmail))
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
