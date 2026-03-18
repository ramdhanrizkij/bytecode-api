package service

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/ramdhanrizki/bytecode-api/internal/identity/application/dto"
	identityDomain "github.com/ramdhanrizki/bytecode-api/internal/identity/domain"
	"github.com/ramdhanrizki/bytecode-api/internal/identity/domain/repository"
	sharedErrors "github.com/ramdhanrizki/bytecode-api/internal/shared/errors"
)

type ProfileService struct {
	users repository.UserRepository
}

func NewProfileService(users repository.UserRepository) *ProfileService {
	return &ProfileService{users: users}
}

func (s *ProfileService) GetCurrent(ctx context.Context, userID string) (*dto.UserSummary, error) {
	parsedID, err := uuid.Parse(userID)
	if err != nil {
		return nil, sharedErrors.Validation("invalid user id", map[string][]string{"user_id": {"user id must be a valid uuid"}})
	}

	user, err := s.users.FindByIDWithRoles(ctx, parsedID)
	if err != nil {
		return nil, sharedErrors.Internal("failed to load profile", err)
	}
	if user == nil {
		return nil, identityDomain.ErrUserNotFound
	}

	summary := toUserSummary(*user)
	return &summary, nil
}

func (s *ProfileService) UpdateCurrent(ctx context.Context, input dto.UpdateProfileInput) (*dto.UserSummary, error) {
	parsedID, err := uuid.Parse(input.UserID)
	if err != nil {
		return nil, sharedErrors.Validation("invalid user id", map[string][]string{"user_id": {"user id must be a valid uuid"}})
	}

	user, err := s.users.FindByIDWithRoles(ctx, parsedID)
	if err != nil {
		return nil, sharedErrors.Internal("failed to load profile", err)
	}
	if user == nil {
		return nil, identityDomain.ErrUserNotFound
	}

	user.FullName = strings.TrimSpace(input.FullName)
	user.UpdatedAt = time.Now().UTC()
	if err := s.users.Update(ctx, user); err != nil {
		return nil, sharedErrors.Internal("failed to update profile", err)
	}

	summary := toUserSummary(*user)
	return &summary, nil
}
