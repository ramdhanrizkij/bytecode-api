package service

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/ramdhanrizkij/bytecode-api/internal/core/worker"
	"github.com/ramdhanrizkij/bytecode-api/internal/features/auth/domain"
	"github.com/ramdhanrizkij/bytecode-api/internal/features/auth/job"
	"github.com/ramdhanrizkij/bytecode-api/internal/model"
	apperrors "github.com/ramdhanrizkij/bytecode-api/internal/shared/errors"
	"github.com/ramdhanrizkij/bytecode-api/pkg/hash"
	pkgjwt "github.com/ramdhanrizkij/bytecode-api/pkg/jwt"
)

// authService implements domain.AuthService with transport-agnostic business logic.
type authService struct {
	repo            domain.AuthRepository
	workerPool      *worker.WorkerPool
	jwtSecret       string
	jwtExpHours     int
	refreshExpHours int
	log             *zap.Logger
}

// NewAuthService creates an AuthService with the required dependencies injected.
func NewAuthService(
	repo domain.AuthRepository,
	wp *worker.WorkerPool,
	jwtSecret string,
	jwtExpHours int,
	refreshExpHours int,
	log *zap.Logger,
) domain.AuthService {
	return &authService{
		repo:            repo,
		workerPool:      wp,
		jwtSecret:       jwtSecret,
		jwtExpHours:     jwtExpHours,
		refreshExpHours: refreshExpHours,
		log:             log,
	}
}

// Register creates a new user account with the default "user" role.
func (s *authService) Register(ctx context.Context, req *domain.RegisterRequest) (*domain.AuthResponse, error) {
	// 1. Check for duplicate email.
	existing, err := s.repo.FindUserByEmail(ctx, req.Email)
	if err != nil && !errors.Is(err, apperrors.ErrNotFound) {
		return nil, apperrors.WrapError(err, "failed to check existing user")
	}
	if existing != nil {
		return nil, apperrors.NewAppError(409, "email already registered", nil)
	}

	// 2. Resolve default role.
	role, err := s.repo.FindRoleByName(ctx, "user")
	if err != nil {
		s.log.Error("default role 'user' not found", zap.Error(err))
		return nil, apperrors.WrapError(err, "default role not found")
	}

	// 3. Hash password.
	hashedPwd, err := hash.HashPassword(req.Password)
	if err != nil {
		return nil, apperrors.WrapError(err, "failed to hash password")
	}

	// 4. Persist the new user.
	roleID := role.ID
	user := &model.User{
		Name:           req.Name,
		Email:          req.Email,
		Password:       hashedPwd,
		ProfilePicture: normalizeOptionalString(req.ProfilePicture),
		RoleID:         &roleID,
		IsActive:       true,
	}
	if err := s.repo.CreateUser(ctx, user); err != nil {
		return nil, err
	}

	// 5. Submit welcome email job to worker pool.
	emailJob := job.NewSendWelcomeEmailJob(user.Email, user.Name, s.log)
	if err := s.workerPool.Submit(emailJob); err != nil {
		s.log.Error("failed to submit welcome email job", zap.Error(err))
		// We don't return error here because the user is already registered.
	}

	// 6. Generate JWT and refresh token.
	token, err := pkgjwt.GenerateToken(
		user.ID.String(), user.Email, role.Name,
		s.jwtSecret, s.jwtExpHours,
	)
	if err != nil {
		return nil, apperrors.WrapError(err, "failed to generate token")
	}
	refreshToken, err := s.createRefreshToken(ctx, user.ID.String())
	if err != nil {
		return nil, err
	}

	return &domain.AuthResponse{
		User: domain.UserResponse{
			ID:       user.ID.String(),
			Name:     user.Name,
			Email:    user.Email,
			RoleName: role.Name,
			IsActive: user.IsActive,
		},
		Token:        token,
		RefreshToken: refreshToken,
	}, nil
}

// Login authenticates a user and returns a JWT token.
func (s *authService) Login(ctx context.Context, req *domain.LoginRequest) (*domain.AuthResponse, error) {
	// 1. Find user by email.
	user, err := s.repo.FindUserByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
			return nil, apperrors.ErrUnauthorized
		}
		return nil, apperrors.WrapError(err, "failed to find user")
	}

	// 2. Verify password.
	if !hash.CheckPassword(req.Password, user.Password) {
		return nil, apperrors.ErrUnauthorized
	}

	// 3. Check that account is active.
	if !user.IsActive {
		return nil, apperrors.ErrForbidden
	}

	// 4. Resolve role name (Role may be nil if FK is NULL).
	roleName := ""
	if user.Role != nil {
		roleName = user.Role.Name
	}

	// 5. Generate JWT and refresh token.
	token, err := pkgjwt.GenerateToken(
		user.ID.String(), user.Email, roleName,
		s.jwtSecret, s.jwtExpHours,
	)
	if err != nil {
		return nil, apperrors.WrapError(err, "failed to generate token")
	}
	refreshToken, err := s.createRefreshToken(ctx, user.ID.String())
	if err != nil {
		return nil, err
	}

	return &domain.AuthResponse{
		User: domain.UserResponse{
			ID:       user.ID.String(),
			Name:     user.Name,
			Email:    user.Email,
			RoleName: roleName,
			IsActive: user.IsActive,
		},
		Token:        token,
		RefreshToken: refreshToken,
	}, nil
}

// Refresh exchanges a valid refresh token for a new access token and rotated refresh token.
func (s *authService) Refresh(ctx context.Context, req *domain.RefreshTokenRequest) (*domain.TokenResponse, error) {
	storedToken, err := s.repo.FindRefreshTokenByHash(ctx, hashRefreshToken(req.RefreshToken))
	if err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
			return nil, apperrors.ErrUnauthorized
		}
		return nil, apperrors.WrapError(err, "failed to find refresh token")
	}

	if storedToken.RevokedAt != nil || time.Now().After(storedToken.ExpiresAt) {
		return nil, apperrors.ErrUnauthorized
	}
	if storedToken.User == nil || !storedToken.User.IsActive {
		return nil, apperrors.ErrUnauthorized
	}

	if err := s.repo.RevokeRefreshToken(ctx, storedToken.ID.String()); err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
			return nil, apperrors.ErrUnauthorized
		}
		return nil, err
	}

	roleName := ""
	if storedToken.User.Role != nil {
		roleName = storedToken.User.Role.Name
	}

	token, err := pkgjwt.GenerateToken(
		storedToken.User.ID.String(), storedToken.User.Email, roleName,
		s.jwtSecret, s.jwtExpHours,
	)
	if err != nil {
		return nil, apperrors.WrapError(err, "failed to generate token")
	}

	refreshToken, err := s.createRefreshToken(ctx, storedToken.User.ID.String())
	if err != nil {
		return nil, err
	}

	return &domain.TokenResponse{
		Token:        token,
		RefreshToken: refreshToken,
	}, nil
}

// Logout revokes a refresh token so it can no longer be exchanged.
func (s *authService) Logout(ctx context.Context, req *domain.LogoutRequest) error {
	storedToken, err := s.repo.FindRefreshTokenByHash(ctx, hashRefreshToken(req.RefreshToken))
	if err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
			return nil
		}
		return apperrors.WrapError(err, "failed to find refresh token")
	}

	if err := s.repo.RevokeRefreshToken(ctx, storedToken.ID.String()); err != nil && !errors.Is(err, apperrors.ErrNotFound) {
		return err
	}

	return nil
}

func (s *authService) CleanupExpiredTokens(ctx context.Context) error {
	s.log.Info("Cleaning up expired tokens...")
	return s.repo.CleanupExpiredTokens(ctx)
}

func (s *authService) createRefreshToken(ctx context.Context, userID string) (string, error) {
	rawToken, err := generateRefreshToken()
	if err != nil {
		return "", apperrors.WrapError(err, "failed to generate refresh token")
	}
	parsedUserID, err := uuid.Parse(userID)
	if err != nil {
		return "", apperrors.WrapError(err, "invalid user id")
	}

	refreshExpiryHours := s.refreshExpHours
	if refreshExpiryHours <= 0 {
		refreshExpiryHours = 168
	}

	token := &model.RefreshToken{
		UserID:    parsedUserID,
		TokenHash: hashRefreshToken(rawToken),
		ExpiresAt: time.Now().Add(time.Duration(refreshExpiryHours) * time.Hour),
	}
	if err := s.repo.CreateRefreshToken(ctx, token); err != nil {
		return "", err
	}

	return rawToken, nil
}

func generateRefreshToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func hashRefreshToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}

func normalizeOptionalString(value *string) *string {
	if value == nil {
		return nil
	}

	trimmed := *value
	if trimmed == "" {
		return nil
	}

	return &trimmed
}
