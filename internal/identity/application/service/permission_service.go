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
	sharedErrors "github.com/ramdhanrizki/bytecode-api/internal/shared/errors"
)

type PermissionService struct {
	permissions repository.PermissionRepository
}

func NewPermissionService(permissions repository.PermissionRepository) *PermissionService {
	return &PermissionService{permissions: permissions}
}

func (s *PermissionService) List(ctx context.Context, input dto.ListInput) (*dto.PermissionListOutput, error) {
	permissions, total, err := s.permissions.List(ctx, toListOptions(input))
	if err != nil {
		return nil, sharedErrors.Internal("failed to load permissions", err)
	}

	items := make([]dto.PermissionSummary, 0, len(permissions))
	for _, permission := range permissions {
		items = append(items, toPermissionSummary(permission))
	}

	return &dto.PermissionListOutput{
		Permissions: items,
		Meta:        toPaginationMeta(input, total),
	}, nil
}

func (s *PermissionService) Get(ctx context.Context, id string) (*dto.PermissionSummary, error) {
	parsedID, err := parseUUID(id, "id")
	if err != nil {
		return nil, err
	}

	permission, err := s.permissions.FindByID(ctx, parsedID)
	if err != nil {
		return nil, sharedErrors.Internal("failed to load permission", err)
	}
	if permission == nil {
		return nil, identityDomain.ErrPermissionNotFound
	}

	summary := toPermissionSummary(*permission)
	return &summary, nil
}

func (s *PermissionService) Create(ctx context.Context, input dto.CreatePermissionInput) (*dto.PermissionSummary, error) {
	name := normalizeName(input.Name)
	existing, err := s.permissions.FindByName(ctx, name)
	if err != nil {
		return nil, sharedErrors.Internal("failed to check existing permission", err)
	}
	if existing != nil {
		return nil, sharedErrors.Conflict("permission name already exists")
	}

	now := time.Now().UTC()
	permission := &entity.Permission{
		ID:          uuid.New(),
		Name:        name,
		Description: strings.TrimSpace(input.Description),
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := s.permissions.Create(ctx, permission); err != nil {
		return nil, sharedErrors.Internal("failed to create permission", err)
	}

	summary := toPermissionSummary(*permission)
	return &summary, nil
}

func (s *PermissionService) Update(ctx context.Context, input dto.UpdatePermissionInput) (*dto.PermissionSummary, error) {
	parsedID, err := parseUUID(input.ID, "id")
	if err != nil {
		return nil, err
	}

	permission, err := s.permissions.FindByID(ctx, parsedID)
	if err != nil {
		return nil, sharedErrors.Internal("failed to load permission", err)
	}
	if permission == nil {
		return nil, identityDomain.ErrPermissionNotFound
	}

	name := normalizeName(input.Name)
	if !strings.EqualFold(permission.Name, name) {
		existing, err := s.permissions.FindByName(ctx, name)
		if err != nil {
			return nil, sharedErrors.Internal("failed to check existing permission", err)
		}
		if existing != nil && existing.ID != permission.ID {
			return nil, sharedErrors.Conflict("permission name already exists")
		}
	}

	permission.Name = name
	permission.Description = strings.TrimSpace(input.Description)
	permission.UpdatedAt = time.Now().UTC()

	if err := s.permissions.Update(ctx, permission); err != nil {
		return nil, sharedErrors.Internal("failed to update permission", err)
	}

	summary := toPermissionSummary(*permission)
	return &summary, nil
}

func (s *PermissionService) Delete(ctx context.Context, id string) error {
	parsedID, err := parseUUID(id, "id")
	if err != nil {
		return err
	}

	permission, err := s.permissions.FindByID(ctx, parsedID)
	if err != nil {
		return sharedErrors.Internal("failed to load permission", err)
	}
	if permission == nil {
		return identityDomain.ErrPermissionNotFound
	}

	if err := s.permissions.Delete(ctx, parsedID); err != nil {
		return sharedErrors.Internal("failed to delete permission", err)
	}

	return nil
}

func toPermissionSummary(permission entity.Permission) dto.PermissionSummary {
	return dto.PermissionSummary{
		ID:          permission.ID.String(),
		Name:        permission.Name,
		Description: permission.Description,
	}
}
