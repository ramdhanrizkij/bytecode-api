package repository

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"

	"github.com/ramdhanrizkij/bytecode-api/internal/features/auth/domain"
	"github.com/ramdhanrizkij/bytecode-api/internal/model"
	apperrors "github.com/ramdhanrizkij/bytecode-api/internal/shared/errors"
)

// authRepository implements domain.AuthRepository using GORM.
type authRepository struct {
	db *gorm.DB
}

// NewAuthRepository creates a new auth repository backed by the given GORM DB.
func NewAuthRepository(db *gorm.DB) domain.AuthRepository {
	return &authRepository{db: db}
}

// FindUserByEmail retrieves a user by email address, preloading the associated Role.
// Returns ErrNotFound if no matching user exists.
func (r *authRepository) FindUserByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	result := r.db.WithContext(ctx).
		Preload("Role").
		Where("email = ?", email).
		First(&user)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrNotFound
		}
		return nil, apperrors.WrapError(result.Error, "failed to find user by email")
	}

	return &user, nil
}

// CreateUser inserts a new user record into the database.
func (r *authRepository) CreateUser(ctx context.Context, user *model.User) error {
	if err := r.db.WithContext(ctx).Create(user).Error; err != nil {
		return apperrors.WrapError(err, "failed to create user")
	}
	return nil
}

// FindRoleByName retrieves a role by its name.
// Returns ErrNotFound if no matching role exists.
func (r *authRepository) FindRoleByName(ctx context.Context, name string) (*model.Role, error) {
	var role model.Role
	result := r.db.WithContext(ctx).
		Where("name = ?", name).
		First(&role)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrNotFound
		}
		return nil, apperrors.WrapError(result.Error, "failed to find role by name")
	}

	return &role, nil
}

// CreateRefreshToken inserts a hashed refresh token record.
func (r *authRepository) CreateRefreshToken(ctx context.Context, token *model.RefreshToken) error {
	if err := r.db.WithContext(ctx).Create(token).Error; err != nil {
		return apperrors.WrapError(err, "failed to create refresh token")
	}
	return nil
}

// FindRefreshTokenByHash retrieves a refresh token record by its hashed token value.
func (r *authRepository) FindRefreshTokenByHash(ctx context.Context, tokenHash string) (*model.RefreshToken, error) {
	var token model.RefreshToken
	result := r.db.WithContext(ctx).
		Preload("User").
		Preload("User.Role").
		Where("token_hash = ?", tokenHash).
		First(&token)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrNotFound
		}
		return nil, apperrors.WrapError(result.Error, "failed to find refresh token")
	}

	return &token, nil
}

// RevokeRefreshToken marks a refresh token as revoked.
func (r *authRepository) RevokeRefreshToken(ctx context.Context, id string) error {
	now := time.Now()
	result := r.db.WithContext(ctx).
		Model(&model.RefreshToken{}).
		Where("id = ? AND revoked_at IS NULL", id).
		Update("revoked_at", now)

	if result.Error != nil {
		return apperrors.WrapError(result.Error, "failed to revoke refresh token")
	}
	if result.RowsAffected == 0 {
		return apperrors.ErrNotFound
	}

	return nil
}

func (r *authRepository) CleanupExpiredTokens(ctx context.Context) error {
	if err := r.db.WithContext(ctx).
		Where("expires_at < ? OR revoked_at IS NOT NULL", time.Now()).
		Delete(&model.RefreshToken{}).Error; err != nil {
		return apperrors.WrapError(err, "failed to cleanup expired refresh tokens")
	}
	return nil
}
