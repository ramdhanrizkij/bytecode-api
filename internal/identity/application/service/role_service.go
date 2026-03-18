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

type RoleService struct {
	roles       repository.RoleRepository
	permissions repository.PermissionRepository
	unitOfWork  repository.UnitOfWork
}

func NewRoleService(roles repository.RoleRepository, permissions repository.PermissionRepository, unitOfWork repository.UnitOfWork) *RoleService {
	return &RoleService{
		roles:       roles,
		permissions: permissions,
		unitOfWork:  unitOfWork,
	}
}

func (s *RoleService) List(ctx context.Context, input dto.ListInput) (*dto.RoleListOutput, error) {
	roles, total, err := s.roles.List(ctx, toListOptions(input))
	if err != nil {
		return nil, sharedErrors.Internal("failed to load roles", err)
	}

	items := make([]dto.RoleSummary, 0, len(roles))
	for _, role := range roles {
		items = append(items, toRoleSummary(role))
	}

	return &dto.RoleListOutput{
		Roles: items,
		Meta:  toPaginationMeta(input, total),
	}, nil
}

func (s *RoleService) Get(ctx context.Context, id string) (*dto.RoleSummary, error) {
	parsedID, err := parseUUID(id, "id")
	if err != nil {
		return nil, err
	}

	role, err := s.roles.FindByIDWithPermissions(ctx, parsedID)
	if err != nil {
		return nil, sharedErrors.Internal("failed to load role", err)
	}
	if role == nil {
		return nil, identityDomain.ErrRoleNotFound
	}

	summary := toRoleSummary(*role)
	return &summary, nil
}

func (s *RoleService) Create(ctx context.Context, input dto.CreateRoleInput) (*dto.RoleSummary, error) {
	name := normalizeName(input.Name)
	existing, err := s.roles.FindByName(ctx, name)
	if err != nil {
		return nil, sharedErrors.Internal("failed to check existing role", err)
	}
	if existing != nil {
		return nil, sharedErrors.Conflict("role name already exists")
	}

	now := time.Now().UTC()
	role := &entity.Role{
		ID:          uuid.New(),
		Name:        name,
		Description: strings.TrimSpace(input.Description),
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := s.roles.Create(ctx, role); err != nil {
		return nil, sharedErrors.Internal("failed to create role", err)
	}

	summary := toRoleSummary(*role)
	return &summary, nil
}

func (s *RoleService) Update(ctx context.Context, input dto.UpdateRoleInput) (*dto.RoleSummary, error) {
	parsedID, err := parseUUID(input.ID, "id")
	if err != nil {
		return nil, err
	}

	role, err := s.roles.FindByIDWithPermissions(ctx, parsedID)
	if err != nil {
		return nil, sharedErrors.Internal("failed to load role", err)
	}
	if role == nil {
		return nil, identityDomain.ErrRoleNotFound
	}

	name := normalizeName(input.Name)
	if !strings.EqualFold(role.Name, name) {
		existing, err := s.roles.FindByName(ctx, name)
		if err != nil {
			return nil, sharedErrors.Internal("failed to check existing role", err)
		}
		if existing != nil && existing.ID != role.ID {
			return nil, sharedErrors.Conflict("role name already exists")
		}
	}

	role.Name = name
	role.Description = strings.TrimSpace(input.Description)
	role.UpdatedAt = time.Now().UTC()

	if err := s.roles.Update(ctx, role); err != nil {
		return nil, sharedErrors.Internal("failed to update role", err)
	}

	summary := toRoleSummary(*role)
	return &summary, nil
}

func (s *RoleService) Delete(ctx context.Context, id string) error {
	parsedID, err := parseUUID(id, "id")
	if err != nil {
		return err
	}

	role, err := s.roles.FindByID(ctx, parsedID)
	if err != nil {
		return sharedErrors.Internal("failed to load role", err)
	}
	if role == nil {
		return identityDomain.ErrRoleNotFound
	}

	if err := s.roles.Delete(ctx, parsedID); err != nil {
		return sharedErrors.Internal("failed to delete role", err)
	}

	return nil
}

func (s *RoleService) AssignPermissions(ctx context.Context, input dto.AssignRolePermissionsInput) (*dto.RoleSummary, error) {
	parsedRoleID, err := parseUUID(input.RoleID, "role_id")
	if err != nil {
		return nil, err
	}

	role, err := s.roles.FindByIDWithPermissions(ctx, parsedRoleID)
	if err != nil {
		return nil, sharedErrors.Internal("failed to load role", err)
	}
	if role == nil {
		return nil, identityDomain.ErrRoleNotFound
	}

	permissionIDs, err := parseUUIDs(input.PermissionIDs, "permission_ids")
	if err != nil {
		return nil, err
	}

	permissions := make([]entity.Permission, 0, len(permissionIDs))
	if len(permissionIDs) > 0 {
		permissions, err = s.permissions.FindByIDs(ctx, permissionIDs)
		if err != nil {
			return nil, sharedErrors.Internal("failed to load permissions", err)
		}
		if len(permissions) != len(permissionIDs) {
			return nil, identityDomain.ErrPermissionNotFound
		}
	}

	if err := s.unitOfWork.Do(ctx, func(repos repository.TxRepositories) error {
		if err := repos.RolePermissions().ReplacePermissions(ctx, parsedRoleID, permissionIDs); err != nil {
			return sharedErrors.Internal("failed to assign permissions", err)
		}
		return nil
	}); err != nil {
		return nil, err
	}

	role.Permissions = permissions
	summary := toRoleSummary(*role)
	return &summary, nil
}

func toRoleSummary(role entity.Role) dto.RoleSummary {
	permissions := make([]dto.PermissionSummary, 0, len(role.Permissions))
	for _, permission := range role.Permissions {
		permissions = append(permissions, toPermissionSummary(permission))
	}

	return dto.RoleSummary{
		ID:          role.ID.String(),
		Name:        role.Name,
		Description: role.Description,
		Permissions: permissions,
	}
}

func normalizeName(value string) string {
	return strings.ToLower(strings.TrimSpace(value))
}
