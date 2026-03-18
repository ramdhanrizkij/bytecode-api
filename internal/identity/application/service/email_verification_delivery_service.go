package service

import (
	"context"
	"fmt"
	"strings"

	"go.uber.org/zap"

	"github.com/ramdhanrizki/bytecode-api/internal/identity/application/job"
	sharedLogger "github.com/ramdhanrizki/bytecode-api/internal/shared/logger"
)

type EmailVerificationDeliveryService struct {
	logger  sharedLogger.Logger
	sender  MailSender
	appName string
}

func NewEmailVerificationDeliveryService(logger sharedLogger.Logger, sender MailSender, appName string) *EmailVerificationDeliveryService {
	return &EmailVerificationDeliveryService{
		logger:  logger,
		sender:  sender,
		appName: strings.TrimSpace(appName),
	}
}

func (s *EmailVerificationDeliveryService) Send(ctx context.Context, payload job.EmailVerificationJob) error {
	fullName := strings.TrimSpace(payload.FullName)
	if fullName == "" {
		fullName = "there"
	}

	appName := s.appName
	if appName == "" {
		appName = "our service"
	}

	message := MailMessage{
		ToAddress: strings.TrimSpace(payload.Email),
		ToName:    fullName,
		Subject:   fmt.Sprintf("Verify your email for %s", appName),
		Body: fmt.Sprintf("Hello %s,\n\nWelcome to %s. Please verify your email address by opening the link below:\n\n%s\n\nIf you did not create this account, you can ignore this email.\n",
			fullName,
			appName,
			strings.TrimSpace(payload.VerificationURL),
		),
	}

	if err := s.sender.Send(ctx, message); err != nil {
		s.logger.Error("failed to send verification email", zap.String("email", payload.Email), zap.String("user_id", payload.UserID), zap.Error(err))
		return err
	}

	s.logger.Info("verification email sent", zap.String("email", payload.Email), zap.String("user_id", payload.UserID))
	return nil
}
