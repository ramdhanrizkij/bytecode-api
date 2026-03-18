package service

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/ramdhanrizki/bytecode-api/internal/identity/application/dto"
	identityDomain "github.com/ramdhanrizki/bytecode-api/internal/identity/domain"
	"github.com/ramdhanrizki/bytecode-api/internal/identity/domain/entity"
	"github.com/ramdhanrizki/bytecode-api/internal/identity/domain/repository"
	domainService "github.com/ramdhanrizki/bytecode-api/internal/identity/domain/service"
	sharedErrors "github.com/ramdhanrizki/bytecode-api/internal/shared/errors"
	sharedKernel "github.com/ramdhanrizki/bytecode-api/internal/shared/kernel"
)

type AdminUserService struct {
	users          repository.UserRepository
	roles          repository.RoleRepository
	unitOfWork     repository.UnitOfWork
	passwordHasher domainService.PasswordHasher
}

func NewAdminUserService(users repository.UserRepository, roles repository.RoleRepository, unitOfWork repository.UnitOfWork, passwordHasher domainService.PasswordHasher) *AdminUserService {
	return &AdminUserService{
		users:          users,
		roles:          roles,
		unitOfWork:     unitOfWork,
		passwordHasher: passwordHasher,
	}
}

func (s *AdminUserService) List(ctx context.Context, input dto.ListInput) (*dto.UserListOutput, error) {
	users, total, err := s.users.List(ctx, toListOptions(input))
	if err != nil {
		return nil, sharedErrors.Internal("failed to load users", err)
	}

	items := make([]dto.UserSummary, 0, len(users))
	for _, user := range users {
		items = append(items, toUserSummary(user))
	}

	return &dto.UserListOutput{
		Users: items,
		Meta:  toPaginationMeta(input, total),
	}, nil
}

func (s *AdminUserService) Get(ctx context.Context, id string) (*dto.UserSummary, error) {
	parsedID, err := parseUUID(id, "id")
	if err != nil {
		return nil, err
	}

	user, err := s.users.FindByIDWithRoles(ctx, parsedID)
	if err != nil {
		return nil, sharedErrors.Internal("failed to load user", err)
	}
	if user == nil {
		return nil, identityDomain.ErrUserNotFound
	}

	summary := toUserSummary(*user)
	return &summary, nil
}

func (s *AdminUserService) Create(ctx context.Context, input dto.CreateAdminUserInput) (*dto.UserSummary, error) {
	email := normalizeEmail(input.Email)
	existing, err := s.users.FindByEmail(ctx, email)
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

	roles, roleIDs, err := s.resolveRoles(ctx, input.RoleIDs, true)
	if err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	isEmailVerified := false
	if input.IsEmailVerified != nil {
		isEmailVerified = *input.IsEmailVerified
	}
	isActive := true
	if input.IsActive != nil {
		isActive = *input.IsActive
	}

	user := &entity.User{
		ID:              uuid.New(),
		FullName:        strings.TrimSpace(input.FullName),
		Email:           email,
		PasswordHash:    passwordHash,
		IsEmailVerified: isEmailVerified,
		IsActive:        isActive,
		CreatedAt:       now,
		UpdatedAt:       now,
		Roles:           roles,
	}

	if err := s.unitOfWork.Do(ctx, func(repos repository.TxRepositories) error {
		if err := repos.Users().Create(ctx, user); err != nil {
			return sharedErrors.Internal("failed to create user", err)
		}
		if err := repos.UserRoles().ReplaceRoles(ctx, user.ID, roleIDs); err != nil {
			return sharedErrors.Internal("failed to assign roles", err)
		}
		return nil
	}); err != nil {
		return nil, err
	}

	summary := toUserSummary(*user)
	return &summary, nil
}

func (s *AdminUserService) Update(ctx context.Context, input dto.UpdateAdminUserInput) (*dto.UserSummary, error) {
	parsedID, err := parseUUID(input.ID, "id")
	if err != nil {
		return nil, err
	}

	user, err := s.users.FindByIDWithRoles(ctx, parsedID)
	if err != nil {
		return nil, sharedErrors.Internal("failed to load user", err)
	}
	if user == nil {
		return nil, identityDomain.ErrUserNotFound
	}

	email := normalizeEmail(input.Email)
	if !strings.EqualFold(user.Email, email) {
		existing, err := s.users.FindByEmail(ctx, email)
		if err != nil {
			return nil, sharedErrors.Internal("failed to check existing user", err)
		}
		if existing != nil && existing.ID != user.ID {
			return nil, identityDomain.ErrEmailAlreadyRegistered
		}
	}

	user.FullName = strings.TrimSpace(input.FullName)
	user.Email = email
	user.IsEmailVerified = input.IsEmailVerified
	user.IsActive = input.IsActive
	user.UpdatedAt = time.Now().UTC()

	if err := s.users.Update(ctx, user); err != nil {
		return nil, sharedErrors.Internal("failed to update user", err)
	}

	summary := toUserSummary(*user)
	return &summary, nil
}

func (s *AdminUserService) Delete(ctx context.Context, id string) error {
	parsedID, err := parseUUID(id, "id")
	if err != nil {
		return err
	}

	user, err := s.users.FindByID(ctx, parsedID)
	if err != nil {
		return sharedErrors.Internal("failed to load user", err)
	}
	if user == nil {
		return identityDomain.ErrUserNotFound
	}

	if err := s.users.Delete(ctx, parsedID); err != nil {
		return sharedErrors.Internal("failed to delete user", err)
	}

	return nil
}

func (s *AdminUserService) AssignRoles(ctx context.Context, input dto.AssignUserRolesInput) (*dto.UserSummary, error) {
	parsedUserID, err := parseUUID(input.UserID, "user_id")
	if err != nil {
		return nil, err
	}

	user, err := s.users.FindByIDWithRoles(ctx, parsedUserID)
	if err != nil {
		return nil, sharedErrors.Internal("failed to load user", err)
	}
	if user == nil {
		return nil, identityDomain.ErrUserNotFound
	}

	roles, roleIDs, err := s.resolveRoles(ctx, input.RoleIDs, false)
	if err != nil {
		return nil, err
	}

	if err := s.unitOfWork.Do(ctx, func(repos repository.TxRepositories) error {
		if err := repos.UserRoles().ReplaceRoles(ctx, parsedUserID, roleIDs); err != nil {
			return sharedErrors.Internal("failed to assign roles", err)
		}
		return nil
	}); err != nil {
		return nil, err
	}

	user.Roles = roles
	summary := toUserSummary(*user)
	return &summary, nil
}

func (s *AdminUserService) resolveRoles(ctx context.Context, rawRoleIDs []string, useDefault bool) ([]entity.Role, []uuid.UUID, error) {
	if len(rawRoleIDs) == 0 && useDefault {
		defaultRole, err := s.roles.FindByName(ctx, "user")
		if err != nil {
			return nil, nil, sharedErrors.Internal("failed to load default role", err)
		}
		if defaultRole == nil {
			return nil, nil, identityDomain.ErrRoleNotFound
		}
		return []entity.Role{*defaultRole}, []uuid.UUID{defaultRole.ID}, nil
	}

	roleIDs, err := parseUUIDs(rawRoleIDs, "role_ids")
	if err != nil {
		return nil, nil, err
	}
	if len(roleIDs) == 0 {
		return nil, nil, sharedErrors.Validation("validation failed", map[string][]string{"role_ids": {"role_ids must contain at least one role"}})
	}

	roles, err := s.roles.FindByIDs(ctx, roleIDs)
	if err != nil {
		return nil, nil, sharedErrors.Internal("failed to load roles", err)
	}
	if len(roles) != len(roleIDs) {
		return nil, nil, identityDomain.ErrRoleNotFound
	}

	return roles, roleIDs, nil
}

func toListOptions(input dto.ListInput) repository.ListOptions {
	return repository.ListOptions{
		Page:   input.Page,
		Limit:  input.Limit,
		Search: input.Search,
		Sort:   input.Sort,
		Order:  input.Order,
	}
}

func toPaginationMeta(input dto.ListInput, total int) dto.PaginationMeta {
	page := input.Page
	if page <= 0 {
		page = sharedKernel.DefaultPage
	}
	limit := input.Limit
	if limit <= 0 {
		limit = sharedKernel.DefaultLimit
	}

	return dto.PaginationMeta{
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: sharedKernel.TotalPages(total, limit),
	}
}

func parseUUID(raw, field string) (uuid.UUID, error) {
	parsed, err := uuid.Parse(strings.TrimSpace(raw))
	if err != nil {
		return uuid.Nil, sharedErrors.Validation("validation failed", map[string][]string{field: {field + " must be a valid uuid"}})
	}
	return parsed, nil
}

func parseUUIDs(values []string, field string) ([]uuid.UUID, error) {
	seen := make(map[uuid.UUID]struct{}, len(values))
	parsedIDs := make([]uuid.UUID, 0, len(values))
	for _, value := range values {
		parsed, err := parseUUID(value, field)
		if err != nil {
			return nil, err
		}
		if _, exists := seen[parsed]; exists {
			continue
		}
		seen[parsed] = struct{}{}
		parsedIDs = append(parsedIDs, parsed)
	}
	return parsedIDs, nil
}

func normalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}
