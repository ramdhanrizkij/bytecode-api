package repository

import (
	"context"

	"github.com/ramdhanrizki/bytecode-api/internal/identity/domain/entity"
)

type EmailVerificationTokenRepository interface {
	Create(ctx context.Context, token *entity.EmailVerificationToken) error
	FindByToken(ctx context.Context, token string) (*entity.EmailVerificationToken, error)
	Update(ctx context.Context, token *entity.EmailVerificationToken) error
}
