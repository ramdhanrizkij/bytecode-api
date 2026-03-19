package service

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"

	categoryDomain "github.com/ramdhanrizki/bytecode-api/internal/category/domain"
	categoryRepo "github.com/ramdhanrizki/bytecode-api/internal/category/domain/repository"
	"github.com/ramdhanrizki/bytecode-api/internal/product/application/dto"
	productDomain "github.com/ramdhanrizki/bytecode-api/internal/product/domain"
	"github.com/ramdhanrizki/bytecode-api/internal/product/domain/entity"
	productRepo "github.com/ramdhanrizki/bytecode-api/internal/product/domain/repository"
	sharedErrors "github.com/ramdhanrizki/bytecode-api/internal/shared/errors"
	sharedKernel "github.com/ramdhanrizki/bytecode-api/internal/shared/kernel"
)

type ProductService struct {
	products   productRepo.ProductRepository
	categories categoryRepo.CategoryRepository
}

func NewProductService(products productRepo.ProductRepository, categories categoryRepo.CategoryRepository) *ProductService {
	return &ProductService{products: products, categories: categories}
}

func (s *ProductService) List(ctx context.Context, input dto.ListInput) (*dto.ProductListOutput, error) {
	var categoryID *uuid.UUID
	if input.CategoryID != nil && strings.TrimSpace(*input.CategoryID) != "" {
		parsed, err := parseUUID(*input.CategoryID, "category_id")
		if err != nil {
			return nil, err
		}
		categoryID = &parsed
	}

	items, total, err := s.products.List(ctx, productRepo.ListOptions{
		Page:       input.Page,
		Limit:      input.Limit,
		Search:     input.Search,
		Sort:       input.Sort,
		Order:      input.Order,
		CategoryID: categoryID,
		IsActive:   input.IsActive,
	})
	if err != nil {
		return nil, sharedErrors.Internal("failed to load products", err)
	}

	result := make([]dto.ProductSummary, 0, len(items))
	for _, product := range items {
		result = append(result, toProductSummary(product))
	}

	return &dto.ProductListOutput{
		Products: result,
		Meta: dto.PaginationMeta{
			Page:       normalizePage(input.Page),
			Limit:      normalizeLimit(input.Limit),
			Total:      total,
			TotalPages: sharedKernel.TotalPages(total, normalizeLimit(input.Limit)),
		},
	}, nil
}

func (s *ProductService) Get(ctx context.Context, id string) (*dto.ProductSummary, error) {
	parsedID, err := parseUUID(id, "id")
	if err != nil {
		return nil, err
	}

	product, err := s.products.FindByID(ctx, parsedID)
	if err != nil {
		return nil, sharedErrors.Internal("failed to load product", err)
	}
	if product == nil {
		return nil, productDomain.ErrProductNotFound
	}

	summary := toProductSummary(*product)
	return &summary, nil
}

func (s *ProductService) Create(ctx context.Context, input dto.CreateProductInput) (*dto.ProductSummary, error) {
	categoryID, categoryName, err := s.ensureCategory(ctx, input.CategoryID)
	if err != nil {
		return nil, err
	}

	slug := generateSlug(input.Slug, input.Name)
	if err := s.ensureUnique(ctx, uuid.Nil, input.SKU, slug); err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	product := &entity.Product{
		ID:           uuid.New(),
		CategoryID:   categoryID,
		CategoryName: categoryName,
		Name:         strings.TrimSpace(input.Name),
		Slug:         slug,
		Description:  sanitizeDescription(input.Description),
		SKU:          strings.TrimSpace(input.SKU),
		Price:        input.Price,
		Stock:        input.Stock,
		IsActive:     input.IsActive,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	if err := s.products.Create(ctx, product); err != nil {
		return nil, sharedErrors.Internal("failed to create product", err)
	}

	summary := toProductSummary(*product)
	return &summary, nil
}

func (s *ProductService) Update(ctx context.Context, input dto.UpdateProductInput) (*dto.ProductSummary, error) {
	parsedID, err := parseUUID(input.ID, "id")
	if err != nil {
		return nil, err
	}

	product, err := s.products.FindByID(ctx, parsedID)
	if err != nil {
		return nil, sharedErrors.Internal("failed to load product", err)
	}
	if product == nil {
		return nil, productDomain.ErrProductNotFound
	}

	categoryID, categoryName, err := s.ensureCategory(ctx, input.CategoryID)
	if err != nil {
		return nil, err
	}

	slug := generateSlug(input.Slug, input.Name)
	if err := s.ensureUnique(ctx, product.ID, input.SKU, slug); err != nil {
		return nil, err
	}

	product.CategoryID = categoryID
	product.CategoryName = categoryName
	product.Name = strings.TrimSpace(input.Name)
	product.Slug = slug
	product.Description = sanitizeDescription(input.Description)
	product.SKU = strings.TrimSpace(input.SKU)
	product.Price = input.Price
	product.Stock = input.Stock
	product.IsActive = input.IsActive
	product.UpdatedAt = time.Now().UTC()

	if err := s.products.Update(ctx, product); err != nil {
		return nil, sharedErrors.Internal("failed to update product", err)
	}

	summary := toProductSummary(*product)
	return &summary, nil
}

func (s *ProductService) Delete(ctx context.Context, id string) error {
	parsedID, err := parseUUID(id, "id")
	if err != nil {
		return err
	}

	product, err := s.products.FindByID(ctx, parsedID)
	if err != nil {
		return sharedErrors.Internal("failed to load product", err)
	}
	if product == nil {
		return productDomain.ErrProductNotFound
	}

	if err := s.products.Delete(ctx, parsedID); err != nil {
		return sharedErrors.Internal("failed to delete product", err)
	}

	return nil
}

func (s *ProductService) ensureCategory(ctx context.Context, rawID string) (uuid.UUID, string, error) {
	categoryID, err := parseUUID(rawID, "category_id")
	if err != nil {
		return uuid.Nil, "", err
	}

	category, err := s.categories.FindByID(ctx, categoryID)
	if err != nil {
		return uuid.Nil, "", sharedErrors.Internal("failed to load category", err)
	}
	if category == nil {
		return uuid.Nil, "", categoryDomain.ErrCategoryNotFound
	}

	return category.ID, category.Name, nil
}

func (s *ProductService) ensureUnique(ctx context.Context, currentID uuid.UUID, sku, slug string) error {
	bySKU, err := s.products.FindBySKU(ctx, strings.TrimSpace(sku))
	if err != nil {
		return sharedErrors.Internal("failed to check product sku", err)
	}
	if bySKU != nil && bySKU.ID != currentID {
		return sharedErrors.Conflict("product sku already exists")
	}

	bySlug, err := s.products.FindBySlug(ctx, slug)
	if err != nil {
		return sharedErrors.Internal("failed to check product slug", err)
	}
	if bySlug != nil && bySlug.ID != currentID {
		return sharedErrors.Conflict("product slug already exists")
	}

	return nil
}

func toProductSummary(product entity.Product) dto.ProductSummary {
	return dto.ProductSummary{
		ID:           product.ID.String(),
		CategoryID:   product.CategoryID.String(),
		CategoryName: product.CategoryName,
		Name:         product.Name,
		Slug:         product.Slug,
		Description:  product.Description,
		SKU:          product.SKU,
		Price:        product.Price,
		Stock:        product.Stock,
		IsActive:     product.IsActive,
		CreatedAt:    product.CreatedAt,
		UpdatedAt:    product.UpdatedAt,
	}
}

func parseUUID(raw, field string) (uuid.UUID, error) {
	parsed, err := uuid.Parse(strings.TrimSpace(raw))
	if err != nil {
		return uuid.Nil, sharedErrors.Validation("validation failed", map[string][]string{field: {field + " must be a valid uuid"}})
	}
	return parsed, nil
}

func normalizePage(page int) int {
	if page <= 0 {
		return sharedKernel.DefaultPage
	}
	return page
}

func normalizeLimit(limit int) int {
	if limit <= 0 {
		return sharedKernel.DefaultLimit
	}
	if limit > sharedKernel.MaxLimit {
		return sharedKernel.MaxLimit
	}
	return limit
}

func sanitizeDescription(description *string) *string {
	if description == nil {
		return nil
	}
	trimmed := strings.TrimSpace(*description)
	if trimmed == "" {
		return nil
	}
	return &trimmed
}

func generateSlug(slug, fallback string) string {
	value := strings.TrimSpace(slug)
	if value == "" {
		value = fallback
	}
	return slugify(value)
}

func slugify(value string) string {
	value = strings.ToLower(strings.TrimSpace(value))
	var builder strings.Builder
	lastDash := false
	for _, char := range value {
		switch {
		case char >= 'a' && char <= 'z', char >= '0' && char <= '9':
			builder.WriteRune(char)
			lastDash = false
		case char == ' ' || char == '-' || char == '_' || char == '/':
			if builder.Len() > 0 && !lastDash {
				builder.WriteByte('-')
				lastDash = true
			}
		}
	}
	result := strings.Trim(builder.String(), "-")
	if result == "" {
		return "item"
	}
	return result
}
