package worker

import (
	"context"
	"encoding/json"
	"fmt"

	"go.uber.org/zap"

	"github.com/ramdhanrizki/bytecode-api/internal/identity/application/job"
	identityService "github.com/ramdhanrizki/bytecode-api/internal/identity/application/service"
	sharedLogger "github.com/ramdhanrizki/bytecode-api/internal/shared/logger"
	sharedQueue "github.com/ramdhanrizki/bytecode-api/internal/shared/queue"
)

type EmailVerificationHandler struct {
	service *identityService.EmailVerificationDeliveryService
	logger  sharedLogger.Logger
}

func NewEmailVerificationHandler(service *identityService.EmailVerificationDeliveryService, logger sharedLogger.Logger) *EmailVerificationHandler {
	return &EmailVerificationHandler{service: service, logger: logger}
}

func (h *EmailVerificationHandler) Handle(ctx context.Context, queueJob sharedQueue.Job) error {
	var payload job.EmailVerificationJob
	if err := json.Unmarshal(queueJob.Payload, &payload); err != nil {
		return fmt.Errorf("decode email verification job: %w", err)
	}

	h.logger.Info("handling email verification job", zap.String("email", payload.Email), zap.String("user_id", payload.UserID), zap.Int("attempt", queueJob.Attempts))
	return h.service.Send(ctx, payload)
}
