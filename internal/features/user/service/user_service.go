package service

import (
	"context"
	"errors"

	"go.uber.org/zap"

	"github.com/ramdhanrizkij/bytecode-api/internal/features/user/domain"
	"github.com/ramdhanrizkij/bytecode-api/internal/model"
	apperrors "github.com/ramdhanrizkij/bytecode-api/internal/shared/errors"
	"github.com/ramdhanrizkij/bytecode-api/internal/shared/pagination"
	"github.com/ramdhanrizkij/bytecode-api/internal/shared/response"
	"github.com/ramdhanrizkij/bytecode-api/pkg/hash"
	"github.com/google/uuid"
)

type userService struct {
	repo domain.UserRepository
	log  *zap.Logger
}

// NewUserService creates a new UserService instance.
func NewUserService(repo domain.UserRepository, log *zap.Logger) domain.UserService {
	return &userService{repo: repo, log: log}
}

func (s *userService) GetAll(ctx context.Context, pq *pagination.PaginationQuery) ([]domain.UserDetailResponse, *response.PaginationMeta, error) {
	users, total, err := s.repo.FindAll(ctx, pq)
	if err != nil {
		return nil, nil, err
	}

	resp := make([]domain.UserDetailResponse, 0, len(users))
	for _, u := range users {
		resp = append(resp, toUserDetailResponse(u))
	}

	meta := &response.PaginationMeta{
		CurrentPage: pq.Page,
		PerPage:     pq.PerPage,
		TotalItems:  total,
		TotalPages:  pagination.CalculateTotalPages(total, pq.PerPage),
	}

	return resp, meta, nil
}

func (s *userService) GetByID(ctx context.Context, id string) (*domain.UserDetailResponse, error) {
	user, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	resp := toUserDetailResponse(*user)
	return &resp, nil
}

func (s *userService) Create(ctx context.Context, req *domain.CreateUserRequest) (*domain.UserDetailResponse, error) {
	// Check email uniqueness.
	existing, err := s.repo.FindByEmail(ctx, req.Email)
	if err != nil && !errors.Is(err, apperrors.ErrNotFound) {
		return nil, apperrors.WrapError(err, "failed to check email uniqueness")
	}
	if existing != nil {
		return nil, apperrors.NewAppError(409, "email already exists", nil)
	}

	// Hash password.
	hashedPwd, err := hash.HashPassword(req.Password)
	if err != nil {
		return nil, apperrors.WrapError(err, "failed to hash password")
	}

	roleID, err := uuid.Parse(req.RoleID)
	if err != nil {
		return nil, apperrors.NewAppError(400, "invalid role_id format", err)
	}

	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	user := &model.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: hashedPwd,
		RoleID:   &roleID,
		IsActive: isActive,
	}

	if err := s.repo.Create(ctx, user); err != nil {
		return nil, err
	}

	// Re-fetch to get preloaded Role info.
	createdUser, err := s.repo.FindByID(ctx, user.ID.String())
	if err != nil {
		return nil, err
	}

	resp := toUserDetailResponse(*createdUser)
	return &resp, nil
}

func (s *userService) Update(ctx context.Context, id string, req *domain.UpdateUserRequest) (*domain.UserDetailResponse, error) {
	user, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Check email uniqueness if email is changed.
	if user.Email != req.Email {
		existing, err := s.repo.FindByEmail(ctx, req.Email)
		if err != nil && !errors.Is(err, apperrors.ErrNotFound) {
			return nil, apperrors.WrapError(err, "failed to check email uniqueness")
		}
		if existing != nil {
			return nil, apperrors.NewAppError(409, "email already exists", nil)
		}
	}

	user.Name = req.Name
	user.Email = req.Email

	if req.Password != "" {
		hashedPwd, err := hash.HashPassword(req.Password)
		if err != nil {
			return nil, apperrors.WrapError(err, "failed to hash password")
		}
		user.Password = hashedPwd
	}

	roleID, err := uuid.Parse(req.RoleID)
	if err != nil {
		return nil, apperrors.NewAppError(400, "invalid role_id format", err)
	}
	user.RoleID = &roleID

	if req.IsActive != nil {
		user.IsActive = *req.IsActive
	}

	if err := s.repo.Update(ctx, user); err != nil {
		return nil, err
	}

	// Re-fetch to get updated Relation info.
	updatedUser, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	resp := toUserDetailResponse(*updatedUser)
	return &resp, nil
}

func (s *userService) Delete(ctx context.Context, currentUserID string, targetID string) error {
	if currentUserID == targetID {
		return apperrors.NewAppError(403, "cannot delete your own account", nil)
	}

	if _, err := s.repo.FindByID(ctx, targetID); err != nil {
		return err
	}

	return s.repo.Delete(ctx, targetID)
}

func (s *userService) GetPermissions(ctx context.Context, userID string) ([]string, error) {
	perms, err := s.repo.GetPermissions(ctx, userID)
	if err != nil {
		return nil, err
	}

	permissionNames := make([]string, 0, len(perms))
	for _, p := range perms {
		permissionNames = append(permissionNames, p.Name)
	}

	return permissionNames, nil
}

func toUserDetailResponse(u model.User) domain.UserDetailResponse {
	roleInfo := domain.RoleInfo{}
	if u.Role != nil {
		roleInfo.ID = u.Role.ID.String()
		roleInfo.Name = u.Role.Name
	}

	return domain.UserDetailResponse{
		ID:        u.ID.String(),
		Name:      u.Name,
		Email:     u.Email,
		IsActive:  u.IsActive,
		Role:      roleInfo,
		CreatedAt: u.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: u.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}
