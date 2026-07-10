package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.uber.org/zap"

	"github.com/ramdhanrizkij/bytecode-api/internal/core/cache"
	"github.com/ramdhanrizkij/bytecode-api/internal/features/permission/domain"
	"github.com/ramdhanrizkij/bytecode-api/internal/model"
	apperrors "github.com/ramdhanrizkij/bytecode-api/internal/shared/errors"
	"github.com/ramdhanrizkij/bytecode-api/internal/shared/pagination"
	"github.com/ramdhanrizkij/bytecode-api/internal/shared/response"
)

type permissionService struct {
	repo     domain.PermissionRepository
	cache    cache.Client
	cacheTTL time.Duration
	log      *zap.Logger
}

// NewPermissionService creates a PermissionService with the repository injected.
func NewPermissionService(repo domain.PermissionRepository, cacheClient cache.Client, cacheTTL time.Duration, log *zap.Logger) domain.PermissionService {
	return &permissionService{repo: repo, cache: cacheClient, cacheTTL: cacheTTL, log: log}
}

func (s *permissionService) GetAll(ctx context.Context, pq *pagination.PaginationQuery) ([]domain.PermissionDetailResponse, *response.PaginationMeta, error) {
	cacheKey := permissionListCacheKey(pq)
	var cached permissionListCacheEntry
	if hit, err := s.cache.Get(ctx, cacheKey, &cached); err == nil && hit {
		return cached.Permissions, cached.Meta, nil
	} else if err != nil {
		s.log.Warn("failed to read permissions from cache", zap.String("key", cacheKey), zap.Error(err))
	}

	perms, total, err := s.repo.FindAll(ctx, pq)
	if err != nil {
		return nil, nil, err
	}

	resp := make([]domain.PermissionDetailResponse, 0, len(perms))
	for _, p := range perms {
		resp = append(resp, toPermissionResponse(p))
	}

	meta := &response.PaginationMeta{
		CurrentPage: pq.Page,
		PerPage:     pq.PerPage,
		TotalItems:  total,
		TotalPages:  pagination.CalculateTotalPages(total, pq.PerPage),
	}

	if err := s.cache.Set(ctx, cacheKey, permissionListCacheEntry{Permissions: resp, Meta: meta}, s.cacheTTL); err != nil {
		s.log.Warn("failed to store permissions in cache", zap.String("key", cacheKey), zap.Error(err))
	}

	return resp, meta, nil
}

func (s *permissionService) GetByID(ctx context.Context, id string) (*domain.PermissionDetailResponse, error) {
	cacheKey := permissionDetailCacheKey(id)
	var cached domain.PermissionDetailResponse
	if hit, err := s.cache.Get(ctx, cacheKey, &cached); err == nil && hit {
		return &cached, nil
	} else if err != nil {
		s.log.Warn("failed to read permission detail from cache", zap.String("key", cacheKey), zap.Error(err))
	}

	perm, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	resp := toPermissionResponse(*perm)

	if err := s.cache.Set(ctx, cacheKey, resp, s.cacheTTL); err != nil {
		s.log.Warn("failed to store permission detail in cache", zap.String("key", cacheKey), zap.Error(err))
	}

	return &resp, nil
}

func (s *permissionService) Create(ctx context.Context, req *domain.CreatePermissionRequest) (*domain.PermissionDetailResponse, error) {
	// Enforce name uniqueness.
	existing, err := s.repo.FindByName(ctx, req.Name)
	if err != nil && !errors.Is(err, apperrors.ErrNotFound) {
		return nil, apperrors.WrapError(err, "failed to check permission name")
	}
	if existing != nil {
		return nil, apperrors.NewAppError(409, "permission name already exists", nil)
	}

	guardName := req.GuardName
	if guardName == "" {
		guardName = "api"
	}

	perm := &model.Permission{
		Name:        req.Name,
		Description: req.Description,
		GuardName:   guardName,
	}
	if err := s.repo.Create(ctx, perm); err != nil {
		return nil, err
	}

	s.invalidatePermissionCache(ctx, false)

	resp := toPermissionResponse(*perm)
	return &resp, nil
}

func (s *permissionService) Update(ctx context.Context, id string, req *domain.UpdatePermissionRequest) (*domain.PermissionDetailResponse, error) {
	perm, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Unique name check (excluding self).
	if perm.Name != req.Name {
		existing, err := s.repo.FindByName(ctx, req.Name)
		if err != nil && !errors.Is(err, apperrors.ErrNotFound) {
			return nil, apperrors.WrapError(err, "failed to check permission name")
		}
		if existing != nil {
			return nil, apperrors.NewAppError(409, "permission name already exists", nil)
		}
	}

	perm.Name = req.Name
	perm.Description = req.Description
	if req.GuardName != "" {
		perm.GuardName = req.GuardName
	}

	if err := s.repo.Update(ctx, perm); err != nil {
		return nil, err
	}

	s.invalidatePermissionCache(ctx, true)

	resp := toPermissionResponse(*perm)
	return &resp, nil
}

func (s *permissionService) Delete(ctx context.Context, id string) error {
	if _, err := s.repo.FindByID(ctx, id); err != nil {
		return err
	}
	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}

	s.invalidatePermissionCache(ctx, true)
	return nil
}

// toPermissionResponse maps a model.Permission to a domain.PermissionDetailResponse.
func toPermissionResponse(p model.Permission) domain.PermissionDetailResponse {
	return domain.PermissionDetailResponse{
		ID:          p.ID.String(),
		Name:        p.Name,
		Description: p.Description,
		GuardName:   p.GuardName,
		CreatedAt:   p.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:   p.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

type permissionListCacheEntry struct {
	Permissions []domain.PermissionDetailResponse `json:"permissions"`
	Meta        *response.PaginationMeta          `json:"meta"`
}

func permissionListCacheKey(pq *pagination.PaginationQuery) string {
	return fmt.Sprintf(
		"permissions:list:page=%d:per_page=%d:sort=%s:order=%s:search=%s",
		pq.Page,
		pq.PerPage,
		pq.Sort,
		pq.Order,
		pq.Search,
	)
}

func permissionDetailCacheKey(id string) string {
	return fmt.Sprintf("permissions:detail:%s", id)
}

func (s *permissionService) invalidatePermissionCache(ctx context.Context, invalidateRoles bool) {
	if err := s.cache.DeleteByPrefix(ctx, "permissions:"); err != nil {
		s.log.Warn("failed to invalidate permission cache", zap.Error(err))
	}

	if invalidateRoles {
		if err := s.cache.DeleteByPrefix(ctx, "roles:"); err != nil {
			s.log.Warn("failed to invalidate role cache after permission mutation", zap.Error(err))
		}
	}
}
