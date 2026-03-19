package gorm

import (
	"context"
	"errors"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/ramdhanrizki/bytecode-api/internal/product/domain/entity"
	productRepo "github.com/ramdhanrizki/bytecode-api/internal/product/domain/repository"
)

type Repositories struct {
	db *gorm.DB
}

func NewRepositories(db *gorm.DB) *Repositories {
	return &Repositories{db: db}
}

func (r *Repositories) Products() productRepo.ProductRepository {
	return &ProductRepository{db: r.db}
}

type ProductRepository struct {
	db *gorm.DB
}

func (r *ProductRepository) Create(ctx context.Context, product *entity.Product) error {
	return r.db.WithContext(ctx).Create(toProductModel(*product)).Error
}

func (r *ProductRepository) Update(ctx context.Context, product *entity.Product) error {
	return r.db.WithContext(ctx).Model(&ProductModel{}).
		Where("id = ?", product.ID.String()).
		Updates(map[string]any{
			"category_id": product.CategoryID.String(),
			"name":        product.Name,
			"slug":        product.Slug,
			"description": product.Description,
			"sku":         product.SKU,
			"price":       product.Price,
			"stock":       product.Stock,
			"is_active":   product.IsActive,
			"updated_at":  product.UpdatedAt,
		}).Error
}

func (r *ProductRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Where("id = ?", id.String()).Delete(&ProductModel{}).Error
}

func (r *ProductRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.Product, error) {
	var model ProductModel
	err := r.db.WithContext(ctx).Preload("Category").Where("id = ?", id.String()).First(&model).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return toProductEntity(model)
}

func (r *ProductRepository) FindBySKU(ctx context.Context, sku string) (*entity.Product, error) {
	var model ProductModel
	err := r.db.WithContext(ctx).Preload("Category").Where("sku = ?", strings.TrimSpace(sku)).First(&model).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return toProductEntity(model)
}

func (r *ProductRepository) FindBySlug(ctx context.Context, slug string) (*entity.Product, error) {
	var model ProductModel
	err := r.db.WithContext(ctx).Preload("Category").Where("slug = ?", strings.ToLower(strings.TrimSpace(slug))).First(&model).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return toProductEntity(model)
}

func (r *ProductRepository) List(ctx context.Context, options productRepo.ListOptions) ([]entity.Product, int, error) {
	query := r.db.WithContext(ctx).Model(&ProductModel{})
	if search := strings.TrimSpace(options.Search); search != "" {
		like := "%" + search + "%"
		query = query.Where("name ILIKE ?", like)
	}
	if options.CategoryID != nil {
		query = query.Where("category_id = ?", options.CategoryID.String())
	}
	if options.IsActive != nil {
		query = query.Where("is_active = ?", *options.IsActive)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var models []ProductModel
	err := applyListOptions(query.Preload("Category"), options).Find(&models).Error
	if err != nil {
		return nil, 0, err
	}

	items := make([]entity.Product, 0, len(models))
	for _, model := range models {
		product, err := toProductEntity(model)
		if err != nil {
			return nil, 0, err
		}
		items = append(items, *product)
	}

	return items, int(total), nil
}

func applyListOptions(query *gorm.DB, options productRepo.ListOptions) *gorm.DB {
	allowedSorts := map[string]string{
		"name":       "name",
		"price":      "price",
		"stock":      "stock",
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
