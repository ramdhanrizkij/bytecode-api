package service

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/ramdhanrizki/bytecode-api/internal/identity/application/dto"
	"github.com/ramdhanrizki/bytecode-api/internal/identity/application/job"
	identityDomain "github.com/ramdhanrizki/bytecode-api/internal/identity/domain"
	"github.com/ramdhanrizki/bytecode-api/internal/identity/domain/entity"
	"github.com/ramdhanrizki/bytecode-api/internal/identity/domain/repository"
	domainService "github.com/ramdhanrizki/bytecode-api/internal/identity/domain/service"
	sharedErrors "github.com/ramdhanrizki/bytecode-api/internal/shared/errors"
	sharedLogger "github.com/ramdhanrizki/bytecode-api/internal/shared/logger"
)

const defaultVerificationTTL = 24 * time.Hour

type AuthService struct {
	logger              sharedLogger.Logger
	users               repository.UserRepository
	roles               repository.RoleRepository
	refreshTokens       repository.RefreshTokenRepository
	verificationTokens  repository.EmailVerificationTokenRepository
	unitOfWork          repository.UnitOfWork
	passwordHasher      domainService.PasswordHasher
	accessTokenProvider domainService.AccessTokenProvider
	tokenGenerator      domainService.TokenGenerator
	jobPublisher        VerificationJobPublisher
	appBaseURL          string
	refreshTokenTTL     time.Duration
}

type AuthServiceDependencies struct {
	Logger              sharedLogger.Logger
	Users               repository.UserRepository
	Roles               repository.RoleRepository
	RefreshTokens       repository.RefreshTokenRepository
	VerificationTokens  repository.EmailVerificationTokenRepository
	UnitOfWork          repository.UnitOfWork
	PasswordHasher      domainService.PasswordHasher
	AccessTokenProvider domainService.AccessTokenProvider
	TokenGenerator      domainService.TokenGenerator
	JobPublisher        VerificationJobPublisher
	AppBaseURL          string
	RefreshTokenTTL     time.Duration
}

func NewAuthService(deps AuthServiceDependencies) *AuthService {
	return &AuthService{
		logger:              deps.Logger,
		users:               deps.Users,
		roles:               deps.Roles,
		refreshTokens:       deps.RefreshTokens,
		verificationTokens:  deps.VerificationTokens,
		unitOfWork:          deps.UnitOfWork,
		passwordHasher:      deps.PasswordHasher,
		accessTokenProvider: deps.AccessTokenProvider,
		tokenGenerator:      deps.TokenGenerator,
		jobPublisher:        deps.JobPublisher,
		appBaseURL:          strings.TrimRight(deps.AppBaseURL, "/"),
		refreshTokenTTL:     deps.RefreshTokenTTL,
	}
}

func (s *AuthService) Register(ctx context.Context, input dto.RegisterInput) (*dto.RegisterOutput, error) {
	existing, err := s.users.FindByEmail(ctx, strings.ToLower(strings.TrimSpace(input.Email)))
	if err != nil {
		return nil, sharedErrors.Internal("failed to check existing user", err)
	}
	if existing != nil {
		return nil, identityDomain.ErrEmailAlreadyRegistered
	}

	passwordHash, err := s.passwordHasher.Hash(input.Password)
	if err != nil {
		return nil, sharedErrors.Internal("failed to secure password", err)
	}

	now := time.Now().UTC()
	user := &entity.User{
		ID:              uuid.New(),
		FullName:        strings.TrimSpace(input.FullName),
		Email:           strings.ToLower(strings.TrimSpace(input.Email)),
		PasswordHash:    passwordHash,
		IsEmailVerified: false,
		IsActive:        true,
		CreatedAt:       now,
		UpdatedAt:       now,
	}

	verificationTokenValue, err := s.tokenGenerator.Generate()
	if err != nil {
		return nil, sharedErrors.Internal("failed to generate verification token", err)
	}

	verificationToken := &entity.EmailVerificationToken{
		ID:        uuid.New(),
		UserID:    user.ID,
		Token:     verificationTokenValue,
		ExpiresAt: now.Add(defaultVerificationTTL),
		CreatedAt: now,
	}

	if err := s.unitOfWork.Do(ctx, func(repos repository.TxRepositories) error {
		role, err := repos.Roles().FindByName(ctx, "user")
		if err != nil {
			return sharedErrors.Internal("failed to load default role", err)
		}
		if role == nil {
			return identityDomain.ErrRoleNotFound
		}

		if err := repos.Users().Create(ctx, user); err != nil {
			return sharedErrors.Internal("failed to create user", err)
		}
		if err := repos.UserRoles().AssignRole(ctx, user.ID, role.ID); err != nil {
			return sharedErrors.Internal("failed to assign default role", err)
		}
		if err := repos.EmailVerificationTokens().Create(ctx, verificationToken); err != nil {
			return sharedErrors.Internal("failed to create verification token", err)
		}
		user.Roles = []entity.Role{*role}
		return nil
	}); err != nil {
		return nil, err
	}

	payload := job.EmailVerificationJob{
		UserID:          user.ID.String(),
		Email:           user.Email,
		FullName:        user.FullName,
		Token:           verificationToken.Token,
		VerificationURL: s.appBaseURL + "/api/v1/auth/verify-email?token=" + verificationToken.Token,
	}

	publishCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.jobPublisher.PublishEmailVerification(publishCtx, payload); err != nil {
		s.logger.Error("failed to publish email verification job",
			zap.String("user_id", user.ID.String()),
			zap.String("email", user.Email),
			zap.Error(err),
		)
	}

	return &dto.RegisterOutput{
		User:                 toUserSummary(*user),
		VerificationQueued:   true,
		VerificationTokenTTL: verificationToken.ExpiresAt,
	}, nil
}

func (s *AuthService) VerifyEmail(ctx context.Context, input dto.VerifyEmailInput) error {
	token, err := s.verificationTokens.FindByToken(ctx, strings.TrimSpace(input.Token))
	if err != nil {
		return sharedErrors.Internal("failed to load verification token", err)
	}
	if token == nil {
		return identityDomain.ErrInvalidVerificationToken
	}

	now := time.Now().UTC()
	if token.IsUsed() {
		user, loadErr := s.users.FindByID(ctx, token.UserID)
		if loadErr == nil && user != nil && user.IsEmailVerified {
			return nil
		}
		return identityDomain.ErrInvalidVerificationToken
	}
	if token.IsExpired(now) {
		return identityDomain.ErrInvalidVerificationToken
	}

	return s.unitOfWork.Do(ctx, func(repos repository.TxRepositories) error {
		user, err := repos.Users().FindByID(ctx, token.UserID)
		if err != nil {
			return sharedErrors.Internal("failed to load user", err)
		}
		if user == nil {
			return identityDomain.ErrUserNotFound
		}
		if user.IsEmailVerified {
			usedAt := now
			token.UsedAt = &usedAt
			if err := repos.EmailVerificationTokens().Update(ctx, token); err != nil {
				return sharedErrors.Internal("failed to update verification token", err)
			}
			return nil
		}

		user.IsEmailVerified = true
		user.UpdatedAt = now
		if err := repos.Users().Update(ctx, user); err != nil {
			return sharedErrors.Internal("failed to update user email status", err)
		}

		usedAt := now
		token.UsedAt = &usedAt
		if err := repos.EmailVerificationTokens().Update(ctx, token); err != nil {
			return sharedErrors.Internal("failed to update verification token", err)
		}

		return nil
	})
}

func (s *AuthService) Login(ctx context.Context, input dto.LoginInput) (*dto.AuthOutput, error) {
	user, err := s.users.FindByEmailWithRoles(ctx, strings.ToLower(strings.TrimSpace(input.Email)))
	if err != nil {
		return nil, sharedErrors.Internal("failed to load user", err)
	}
	if user == nil {
		return nil, identityDomain.ErrInvalidCredentials
	}
	if err := s.passwordHasher.Compare(user.PasswordHash, input.Password); err != nil {
		return nil, identityDomain.ErrInvalidCredentials
	}
	if !user.IsActive {
		return nil, identityDomain.ErrInactiveUser
	}
	if !user.IsEmailVerified {
		return nil, identityDomain.ErrEmailNotVerified
	}

	return s.issueAuthTokens(ctx, user)
}

func (s *AuthService) RefreshToken(ctx context.Context, input dto.RefreshTokenInput) (*dto.AuthOutput, error) {
	storedToken, err := s.refreshTokens.FindByToken(ctx, strings.TrimSpace(input.RefreshToken))
	if err != nil {
		return nil, sharedErrors.Internal("failed to load refresh token", err)
	}
	if storedToken == nil {
		return nil, identityDomain.ErrInvalidRefreshToken
	}

	now := time.Now().UTC()
	if storedToken.IsRevoked() || storedToken.IsExpired(now) {
		return nil, identityDomain.ErrInvalidRefreshToken
	}

	user, err := s.users.FindByIDWithRoles(ctx, storedToken.UserID)
	if err != nil {
		return nil, sharedErrors.Internal("failed to load user", err)
	}
	if user == nil {
		return nil, identityDomain.ErrUserNotFound
	}
	if !user.IsActive {
		return nil, identityDomain.ErrInactiveUser
	}
	if !user.IsEmailVerified {
		return nil, identityDomain.ErrEmailNotVerified
	}

	var newRefreshToken *entity.RefreshToken
	var accessToken domainService.AccessToken
	if err := s.unitOfWork.Do(ctx, func(repos repository.TxRepositories) error {
		revokedAt := now
		storedToken.RevokedAt = &revokedAt
		if err := repos.RefreshTokens().Update(ctx, storedToken); err != nil {
			return sharedErrors.Internal("failed to revoke refresh token", err)
		}

		tokenValue, err := s.tokenGenerator.Generate()
		if err != nil {
			return sharedErrors.Internal("failed to generate refresh token", err)
		}

		newRefreshToken = &entity.RefreshToken{
			ID:        uuid.New(),
			UserID:    user.ID,
			Token:     tokenValue,
			ExpiresAt: now.Add(s.refreshTokenTTL),
			CreatedAt: now,
		}

		if err := repos.RefreshTokens().Create(ctx, newRefreshToken); err != nil {
			return sharedErrors.Internal("failed to persist refresh token", err)
		}

		accessToken, err = s.accessTokenProvider.Generate(domainService.AccessTokenClaims{
			UserID: user.ID,
			Email:  user.Email,
			Roles:  roleNames(user.Roles),
		})
		if err != nil {
			return sharedErrors.Internal("failed to generate access token", err)
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return &dto.AuthOutput{
		User: toUserSummary(*user),
		Tokens: dto.AuthTokens{
			AccessToken:           accessToken.Token,
			RefreshToken:          newRefreshToken.Token,
			TokenType:             "Bearer",
			AccessTokenExpiresAt:  accessToken.ExpiresAt,
			RefreshTokenExpiresAt: newRefreshToken.ExpiresAt,
		},
	}, nil
}

func (s *AuthService) issueAuthTokens(ctx context.Context, user *entity.User) (*dto.AuthOutput, error) {
	now := time.Now().UTC()
	accessToken, err := s.accessTokenProvider.Generate(domainService.AccessTokenClaims{
		UserID: user.ID,
		Email:  user.Email,
		Roles:  roleNames(user.Roles),
	})
	if err != nil {
		return nil, sharedErrors.Internal("failed to generate access token", err)
	}

	refreshTokenValue, err := s.tokenGenerator.Generate()
	if err != nil {
		return nil, sharedErrors.Internal("failed to generate refresh token", err)
	}

	refreshToken := &entity.RefreshToken{
		ID:        uuid.New(),
		UserID:    user.ID,
		Token:     refreshTokenValue,
		ExpiresAt: now.Add(s.refreshTokenTTL),
		CreatedAt: now,
	}

	if err := s.refreshTokens.Create(ctx, refreshToken); err != nil {
		return nil, sharedErrors.Internal("failed to persist refresh token", err)
	}

	return &dto.AuthOutput{
		User: toUserSummary(*user),
		Tokens: dto.AuthTokens{
			AccessToken:           accessToken.Token,
			RefreshToken:          refreshToken.Token,
			TokenType:             "Bearer",
			AccessTokenExpiresAt:  accessToken.ExpiresAt,
			RefreshTokenExpiresAt: refreshToken.ExpiresAt,
		},
	}, nil
}

func roleNames(roles []entity.Role) []string {
	names := make([]string, 0, len(roles))
	for _, role := range roles {
		names = append(names, role.Name)
	}
	return names
}

func toUserSummary(user entity.User) dto.UserSummary {
	return dto.UserSummary{
		ID:              user.ID.String(),
		FullName:        user.FullName,
		Email:           user.Email,
		IsEmailVerified: user.IsEmailVerified,
		IsActive:        user.IsActive,
		Roles:           roleNames(user.Roles),
	}
}
