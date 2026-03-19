package gorm

import (
	"context"
	"errors"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/ramdhanrizki/bytecode-api/internal/category/domain/entity"
	categoryRepo "github.com/ramdhanrizki/bytecode-api/internal/category/domain/repository"
)

type Repositories struct {
	db *gorm.DB
}

func NewRepositories(db *gorm.DB) *Repositories {
	return &Repositories{db: db}
}

func (r *Repositories) Categories() categoryRepo.CategoryRepository {
	return &CategoryRepository{db: r.db}
}

type CategoryRepository struct {
	db *gorm.DB
}

func (r *CategoryRepository) Create(ctx context.Context, category *entity.Category) error {
	return r.db.WithContext(ctx).Create(toCategoryModel(*category)).Error
}

func (r *CategoryRepository) Update(ctx context.Context, category *entity.Category) error {
	return r.db.WithContext(ctx).Model(&CategoryModel{}).
		Where("id = ?", category.ID.String()).
		Updates(map[string]any{
			"name":        category.Name,
			"slug":        category.Slug,
			"description": category.Description,
			"is_active":   category.IsActive,
			"updated_at":  category.UpdatedAt,
		}).Error
}

func (r *CategoryRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Where("id = ?", id.String()).Delete(&CategoryModel{}).Error
}

func (r *CategoryRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.Category, error) {
	var model CategoryModel
	err := r.db.WithContext(ctx).Where("id = ?", id.String()).First(&model).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return toCategoryEntity(model)
}

func (r *CategoryRepository) FindByName(ctx context.Context, name string) (*entity.Category, error) {
	var model CategoryModel
	err := r.db.WithContext(ctx).Where("LOWER(name) = ?", strings.ToLower(strings.TrimSpace(name))).First(&model).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return toCategoryEntity(model)
}

func (r *CategoryRepository) FindBySlug(ctx context.Context, slug string) (*entity.Category, error) {
	var model CategoryModel
	err := r.db.WithContext(ctx).Where("slug = ?", strings.ToLower(strings.TrimSpace(slug))).First(&model).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return toCategoryEntity(model)
}

func (r *CategoryRepository) List(ctx context.Context, options categoryRepo.ListOptions) ([]entity.Category, int, error) {
	query := r.db.WithContext(ctx).Model(&CategoryModel{})
	if search := strings.TrimSpace(options.Search); search != "" {
		like := "%" + search + "%"
		query = query.Where("name ILIKE ?", like)
	}
	if options.IsActive != nil {
		query = query.Where("is_active = ?", *options.IsActive)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var models []CategoryModel
	err := applyListOptions(query, options).Find(&models).Error
	if err != nil {
		return nil, 0, err
	}

	items := make([]entity.Category, 0, len(models))
	for _, model := range models {
		category, err := toCategoryEntity(model)
		if err != nil {
			return nil, 0, err
		}
		items = append(items, *category)
	}

	return items, int(total), nil
}

func applyListOptions(query *gorm.DB, options categoryRepo.ListOptions) *gorm.DB {
	allowedSorts := map[string]string{
		"name":       "name",
		"created_at": "created_at",
		"updated_at": "updated_at",
	}
	sortColumn := "created_at"
	if candidate, ok := allowedSorts[strings.ToLower(strings.TrimSpace(options.Sort))]; ok {
		sortColumn = candidate
	}
	order := "desc"
	if strings.EqualFold(strings.TrimSpace(options.Order), "asc") {
		order = "asc"
	}
	page := options.Page
	if page <= 0 {
		page = 1
	}
	limit := options.Limit
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}
	return query.Order(sortColumn + " " + order).Offset((page - 1) * limit).Limit(limit)
}
