package service

import (
	"context"
	"errors"

	"go.uber.org/zap"

	"github.com/ramdhanrizkij/bytecode-api/internal/features/permission/domain"
	"github.com/ramdhanrizkij/bytecode-api/internal/model"
	apperrors "github.com/ramdhanrizkij/bytecode-api/internal/shared/errors"
	"github.com/ramdhanrizkij/bytecode-api/internal/shared/pagination"
	"github.com/ramdhanrizkij/bytecode-api/internal/shared/response"
)

type permissionService struct {
	repo domain.PermissionRepository
	log  *zap.Logger
}

// NewPermissionService creates a PermissionService with the repository injected.
func NewPermissionService(repo domain.PermissionRepository, log *zap.Logger) domain.PermissionService {
	return &permissionService{repo: repo, log: log}
}

func (s *permissionService) GetAll(ctx context.Context, pq *pagination.PaginationQuery) ([]domain.PermissionDetailResponse, *response.PaginationMeta, error) {
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

	return resp, meta, nil
}

func (s *permissionService) GetByID(ctx context.Context, id string) (*domain.PermissionDetailResponse, error) {
	perm, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	resp := toPermissionResponse(*perm)
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

	resp := toPermissionResponse(*perm)
	return &resp, nil
}

func (s *permissionService) Delete(ctx context.Context, id string) error {
	if _, err := s.repo.FindByID(ctx, id); err != nil {
		return err
	}
	return s.repo.Delete(ctx, id)
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
