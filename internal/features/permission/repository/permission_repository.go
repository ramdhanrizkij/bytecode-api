package repository

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"github.com/ramdhanrizkij/bytecode-api/internal/features/permission/domain"
	"github.com/ramdhanrizkij/bytecode-api/internal/model"
	apperrors "github.com/ramdhanrizkij/bytecode-api/internal/shared/errors"
	"github.com/ramdhanrizkij/bytecode-api/internal/shared/pagination"
)

type permissionRepository struct {
	db *gorm.DB
}

// NewPermissionRepository creates a PermissionRepository backed by GORM.
func NewPermissionRepository(db *gorm.DB) domain.PermissionRepository {
	return &permissionRepository{db: db}
}

func (r *permissionRepository) FindAll(ctx context.Context, pq *pagination.PaginationQuery) ([]model.Permission, int64, error) {
	var permissions []model.Permission
	var total int64

	query := r.db.WithContext(ctx).Model(&model.Permission{})

	if pq.Search != "" {
		query = query.Where("name ILIKE ?", "%"+pq.Search+"%")
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, apperrors.WrapError(err, "failed to count permissions")
	}

	if err := query.
		Order(pq.GetSort()).
		Limit(pq.GetLimit()).
		Offset(pq.GetOffset()).
		Find(&permissions).Error; err != nil {
		return nil, 0, apperrors.WrapError(err, "failed to fetch permissions")
	}

	return permissions, total, nil
}

func (r *permissionRepository) FindByID(ctx context.Context, id string) (*model.Permission, error) {
	var permission model.Permission
	result := r.db.WithContext(ctx).Where("id = ?", id).First(&permission)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrNotFound
		}
		return nil, apperrors.WrapError(result.Error, "failed to find permission")
	}
	return &permission, nil
}

func (r *permissionRepository) FindByName(ctx context.Context, name string) (*model.Permission, error) {
	var permission model.Permission
	result := r.db.WithContext(ctx).Where("name = ?", name).First(&permission)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrNotFound
		}
		return nil, apperrors.WrapError(result.Error, "failed to find permission by name")
	}
	return &permission, nil
}

// FindByIDs retrieves multiple permissions by their IDs using WHERE id IN (...).
func (r *permissionRepository) FindByIDs(ctx context.Context, ids []string) ([]model.Permission, error) {
	var permissions []model.Permission
	if err := r.db.WithContext(ctx).Where("id IN ?", ids).Find(&permissions).Error; err != nil {
		return nil, apperrors.WrapError(err, "failed to find permissions by IDs")
	}
	return permissions, nil
}

func (r *permissionRepository) Create(ctx context.Context, permission *model.Permission) error {
	if err := r.db.WithContext(ctx).Create(permission).Error; err != nil {
		return apperrors.WrapError(err, "failed to create permission")
	}
	return nil
}

func (r *permissionRepository) Update(ctx context.Context, permission *model.Permission) error {
	if err := r.db.WithContext(ctx).Save(permission).Error; err != nil {
		return apperrors.WrapError(err, "failed to update permission")
	}
	return nil
}

func (r *permissionRepository) Delete(ctx context.Context, id string) error {
	result := r.db.WithContext(ctx).Where("id = ?", id).Delete(&model.Permission{})
	if result.Error != nil {
		return apperrors.WrapError(result.Error, "failed to delete permission")
	}
	if result.RowsAffected == 0 {
		return apperrors.ErrNotFound
	}
	return nil
}
