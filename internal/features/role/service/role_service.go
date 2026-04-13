package service

import (
	"context"
	"errors"

	"go.uber.org/zap"

	"github.com/ramdhanrizkij/bytecode-api/internal/features/role/domain"
	"github.com/ramdhanrizkij/bytecode-api/internal/model"
	apperrors "github.com/ramdhanrizkij/bytecode-api/internal/shared/errors"
	"github.com/ramdhanrizkij/bytecode-api/internal/shared/pagination"
	"github.com/ramdhanrizkij/bytecode-api/internal/shared/response"
)

type roleService struct {
	repo domain.RoleRepository
	log  *zap.Logger
}

// NewRoleService creates a RoleService with the repository injected.
func NewRoleService(repo domain.RoleRepository, log *zap.Logger) domain.RoleService {
	return &roleService{repo: repo, log: log}
}

func (s *roleService) GetAll(ctx context.Context, pq *pagination.PaginationQuery) ([]domain.RoleResponse, *response.PaginationMeta, error) {
	roles, total, err := s.repo.FindAll(ctx, pq)
	if err != nil {
		return nil, nil, err
	}

	resp := make([]domain.RoleResponse, 0, len(roles))
	for _, r := range roles {
		resp = append(resp, toRoleResponse(r))
	}

	meta := &response.PaginationMeta{
		CurrentPage: pq.Page,
		PerPage:     pq.PerPage,
		TotalItems:  total,
		TotalPages:  pagination.CalculateTotalPages(total, pq.PerPage),
	}

	return resp, meta, nil
}

func (s *roleService) GetByID(ctx context.Context, id string) (*domain.RoleResponse, error) {
	role, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	resp := toRoleResponse(*role)
	return &resp, nil
}

func (s *roleService) Create(ctx context.Context, req *domain.CreateRoleRequest) (*domain.RoleResponse, error) {
	// Check name uniqueness.
	existing, err := s.repo.FindByName(ctx, req.Name)
	if err != nil && !errors.Is(err, apperrors.ErrNotFound) {
		return nil, apperrors.WrapError(err, "failed to check role name")
	}
	if existing != nil {
		return nil, apperrors.NewAppError(409, "role name already exists", nil)
	}

	guardName := req.GuardName
	if guardName == "" {
		guardName = "api"
	}

	role := &model.Role{
		Name:        req.Name,
		Description: req.Description,
		GuardName:   guardName,
	}
	if err := s.repo.Create(ctx, role); err != nil {
		return nil, err
	}

	resp := toRoleResponse(*role)
	return &resp, nil
}

func (s *roleService) Update(ctx context.Context, id string, req *domain.UpdateRoleRequest) (*domain.RoleResponse, error) {
	// Verify existence.
	role, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Check name uniqueness (excluding self).
	if role.Name != req.Name {
		existing, err := s.repo.FindByName(ctx, req.Name)
		if err != nil && !errors.Is(err, apperrors.ErrNotFound) {
			return nil, apperrors.WrapError(err, "failed to check role name")
		}
		if existing != nil {
			return nil, apperrors.NewAppError(409, "role name already exists", nil)
		}
	}

	role.Name = req.Name
	role.Description = req.Description
	if req.GuardName != "" {
		role.GuardName = req.GuardName
	}

	if err := s.repo.Update(ctx, role); err != nil {
		return nil, err
	}

	resp := toRoleResponse(*role)
	return &resp, nil
}

func (s *roleService) Delete(ctx context.Context, id string) error {
	// Verify role exists before deleting.
	if _, err := s.repo.FindByID(ctx, id); err != nil {
		return err
	}
	return s.repo.Delete(ctx, id)
}

func (s *roleService) AssignPermissions(ctx context.Context, id string, req *domain.AssignPermissionsRequest) error {
	if _, err := s.repo.FindByID(ctx, id); err != nil {
		return err
	}
	return s.repo.AssignPermissions(ctx, id, req.PermissionIDs)
}

func (s *roleService) RemovePermissions(ctx context.Context, id string, req *domain.AssignPermissionsRequest) error {
	if _, err := s.repo.FindByID(ctx, id); err != nil {
		return err
	}
	return s.repo.RemovePermissions(ctx, id, req.PermissionIDs)
}

// toRoleResponse maps a model.Role to a domain.RoleResponse.
func toRoleResponse(r model.Role) domain.RoleResponse {
	perms := make([]domain.PermissionResponse, 0, len(r.Permissions))
	for _, p := range r.Permissions {
		perms = append(perms, domain.PermissionResponse{
			ID:   p.ID.String(),
			Name: p.Name,
		})
	}

	return domain.RoleResponse{
		ID:          r.ID.String(),
		Name:        r.Name,
		Description: r.Description,
		GuardName:   r.GuardName,
		Permissions: perms,
		CreatedAt:   r.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:   r.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}
