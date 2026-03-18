package repository

import (
	"context"

	"github.com/ramdhanrizki/bytecode-api/internal/identity/domain/entity"
)

type RefreshTokenRepository interface {
	Create(ctx context.Context, token *entity.RefreshToken) error
	FindByToken(ctx context.Context, token string) (*entity.RefreshToken, error)
	Update(ctx context.Context, token *entity.RefreshToken) error
}
