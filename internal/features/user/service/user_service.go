package service

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/ramdhanrizkij/bytecode-api/internal/core/cache"
	"github.com/ramdhanrizkij/bytecode-api/internal/core/storage"
	"github.com/ramdhanrizkij/bytecode-api/internal/features/user/domain"
	"github.com/ramdhanrizkij/bytecode-api/internal/model"
	apperrors "github.com/ramdhanrizkij/bytecode-api/internal/shared/errors"
	"github.com/ramdhanrizkij/bytecode-api/internal/shared/pagination"
	"github.com/ramdhanrizkij/bytecode-api/internal/shared/response"
	"github.com/ramdhanrizkij/bytecode-api/pkg/hash"
)

type userService struct {
	repo     domain.UserRepository
	cache    cache.Client
	storage  storage.Provider
	bucket   string
	cacheTTL time.Duration
	log      *zap.Logger
}

// NewUserService creates a new UserService instance.
func NewUserService(repo domain.UserRepository, cacheClient cache.Client, storageProvider storage.Provider, bucket string, cacheTTL time.Duration, log *zap.Logger) domain.UserService {
	return &userService{
		repo:     repo,
		cache:    cacheClient,
		storage:  storageProvider,
		bucket:   bucket,
		cacheTTL: cacheTTL,
		log:      log,
	}
}

func (s *userService) GetAll(ctx context.Context, pq *pagination.PaginationQuery) ([]domain.UserDetailResponse, *response.PaginationMeta, error) {
	cacheKey := userListCacheKey(pq)
	var cached userListCacheEntry
	if hit, err := s.cache.Get(ctx, cacheKey, &cached); err == nil && hit {
		return cached.Users, cached.Meta, nil
	} else if err != nil {
		s.log.Warn("failed to read users from cache", zap.String("key", cacheKey), zap.Error(err))
	}

	users, total, err := s.repo.FindAll(ctx, pq)
	if err != nil {
		return nil, nil, err
	}

	resp := make([]domain.UserDetailResponse, 0, len(users))
	for _, u := range users {
		resp = append(resp, s.toUserDetailResponse(u))
	}

	meta := &response.PaginationMeta{
		CurrentPage: pq.Page,
		PerPage:     pq.PerPage,
		TotalItems:  total,
		TotalPages:  pagination.CalculateTotalPages(total, pq.PerPage),
	}

	if err := s.cache.Set(ctx, cacheKey, userListCacheEntry{Users: resp, Meta: meta}, s.cacheTTL); err != nil {
		s.log.Warn("failed to store users in cache", zap.String("key", cacheKey), zap.Error(err))
	}

	return resp, meta, nil
}

func (s *userService) GetByID(ctx context.Context, id string) (*domain.UserDetailResponse, error) {
	cacheKey := userDetailCacheKey(id)
	var cached domain.UserDetailResponse
	if hit, err := s.cache.Get(ctx, cacheKey, &cached); err == nil && hit {
		return &cached, nil
	} else if err != nil {
		s.log.Warn("failed to read user detail from cache", zap.String("key", cacheKey), zap.Error(err))
	}

	user, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	resp := s.toUserDetailResponse(*user)

	if err := s.cache.Set(ctx, cacheKey, resp, s.cacheTTL); err != nil {
		s.log.Warn("failed to store user detail in cache", zap.String("key", cacheKey), zap.Error(err))
	}

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
		Name:           req.Name,
		Email:          req.Email,
		Password:       hashedPwd,
		ProfilePicture: normalizeOptionalString(req.ProfilePicture),
		RoleID:         &roleID,
		IsActive:       isActive,
	}

	if err := s.repo.Create(ctx, user); err != nil {
		return nil, err
	}

	// Re-fetch to get preloaded Role info.
	createdUser, err := s.repo.FindByID(ctx, user.ID.String())
	if err != nil {
		return nil, err
	}

	s.invalidateUserCache(ctx)

	resp := s.toUserDetailResponse(*createdUser)
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
	user.ProfilePicture = normalizeOptionalString(req.ProfilePicture)

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

	s.invalidateUserCache(ctx)

	resp := s.toUserDetailResponse(*updatedUser)
	return &resp, nil
}

func (s *userService) UpdateProfile(ctx context.Context, id string, req *domain.UpdateProfileRequest) (*domain.UserDetailResponse, error) {
	user, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

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
	user.ProfilePicture = normalizeOptionalString(req.ProfilePicture)

	if err := s.repo.Update(ctx, user); err != nil {
		return nil, err
	}

	updatedUser, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	s.invalidateUserCache(ctx)

	resp := s.toUserDetailResponse(*updatedUser)
	return &resp, nil
}

func (s *userService) Delete(ctx context.Context, currentUserID string, targetID string) error {
	if currentUserID == targetID {
		return apperrors.NewAppError(403, "cannot delete your own account", nil)
	}

	if _, err := s.repo.FindByID(ctx, targetID); err != nil {
		return err
	}

	if err := s.repo.Delete(ctx, targetID); err != nil {
		return err
	}

	s.invalidateUserCache(ctx)
	return nil
}

func (s *userService) GetPermissions(ctx context.Context, userID string) ([]string, error) {
	cacheKey := userPermissionsCacheKey(userID)
	var cached []string
	if hit, err := s.cache.Get(ctx, cacheKey, &cached); err == nil && hit {
		return cached, nil
	} else if err != nil {
		s.log.Warn("failed to read user permissions from cache", zap.String("key", cacheKey), zap.Error(err))
	}

	perms, err := s.repo.GetPermissions(ctx, userID)
	if err != nil {
		return nil, err
	}

	permissionNames := make([]string, 0, len(perms))
	for _, p := range perms {
		permissionNames = append(permissionNames, p.Name)
	}

	if err := s.cache.Set(ctx, cacheKey, permissionNames, s.cacheTTL); err != nil {
		s.log.Warn("failed to store user permissions in cache", zap.String("key", cacheKey), zap.Error(err))
	}

	return permissionNames, nil
}

func (s *userService) toUserDetailResponse(u model.User) domain.UserDetailResponse {
	roleInfo := domain.RoleInfo{}
	if u.Role != nil {
		roleInfo.ID = u.Role.ID.String()
		roleInfo.Name = u.Role.Name
	}

	return domain.UserDetailResponse{
		ID:             u.ID.String(),
		Name:           u.Name,
		Email:          u.Email,
		ProfilePicture: s.toProfilePictureResponse(u.ProfilePicture),
		IsActive:       u.IsActive,
		Role:           roleInfo,
		CreatedAt:      u.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:      u.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

func (s *userService) toProfilePictureResponse(profilePicture *string) *domain.ProfilePictureResponse {
	if profilePicture == nil {
		return nil
	}

	key := strings.TrimSpace(*profilePicture)
	if key == "" {
		return nil
	}

	return &domain.ProfilePictureResponse{
		Bucket: s.bucket,
		Key:    key,
		URL:    s.storage.URL(s.bucket, key),
	}
}

type userListCacheEntry struct {
	Users []domain.UserDetailResponse `json:"users"`
	Meta  *response.PaginationMeta    `json:"meta"`
}

func userListCacheKey(pq *pagination.PaginationQuery) string {
	return fmt.Sprintf(
		"users:list:page=%d:per_page=%d:sort=%s:order=%s:search=%s",
		pq.Page,
		pq.PerPage,
		pq.Sort,
		pq.Order,
		pq.Search,
	)
}

func userDetailCacheKey(id string) string {
	return fmt.Sprintf("users:detail:%s", id)
}

func userPermissionsCacheKey(id string) string {
	return fmt.Sprintf("users:permissions:%s", id)
}

func (s *userService) invalidateUserCache(ctx context.Context) {
	if err := s.cache.DeleteByPrefix(ctx, "users:"); err != nil {
		s.log.Warn("failed to invalidate user cache", zap.Error(err))
	}
}

func normalizeOptionalString(value *string) *string {
	if value == nil {
		return nil
	}

	trimmed := strings.TrimSpace(*value)
	if trimmed == "" {
		return nil
	}

	return &trimmed
}
