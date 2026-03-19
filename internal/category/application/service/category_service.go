package service

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/ramdhanrizki/bytecode-api/internal/category/application/dto"
	categoryDomain "github.com/ramdhanrizki/bytecode-api/internal/category/domain"
	"github.com/ramdhanrizki/bytecode-api/internal/category/domain/entity"
	"github.com/ramdhanrizki/bytecode-api/internal/category/domain/repository"
	sharedErrors "github.com/ramdhanrizki/bytecode-api/internal/shared/errors"
	sharedKernel "github.com/ramdhanrizki/bytecode-api/internal/shared/kernel"
)

type CategoryService struct {
	categories repository.CategoryRepository
}

func NewCategoryService(categories repository.CategoryRepository) *CategoryService {
	return &CategoryService{categories: categories}
}

func (s *CategoryService) List(ctx context.Context, input dto.ListInput) (*dto.CategoryListOutput, error) {
	items, total, err := s.categories.List(ctx, repository.ListOptions{
		Page:     input.Page,
		Limit:    input.Limit,
		Search:   input.Search,
		Sort:     input.Sort,
		Order:    input.Order,
		IsActive: input.IsActive,
	})
	if err != nil {
		return nil, sharedErrors.Internal("failed to load categories", err)
	}

	result := make([]dto.CategorySummary, 0, len(items))
	for _, category := range items {
		result = append(result, toCategorySummary(category))
	}

	return &dto.CategoryListOutput{
		Categories: result,
		Meta: dto.PaginationMeta{
			Page:       normalizePage(input.Page),
			Limit:      normalizeLimit(input.Limit),
			Total:      total,
			TotalPages: sharedKernel.TotalPages(total, normalizeLimit(input.Limit)),
		},
	}, nil
}

func (s *CategoryService) Get(ctx context.Context, id string) (*dto.CategorySummary, error) {
	parsedID, err := parseUUID(id, "id")
	if err != nil {
		return nil, err
	}

	category, err := s.categories.FindByID(ctx, parsedID)
	if err != nil {
		return nil, sharedErrors.Internal("failed to load category", err)
	}
	if category == nil {
		return nil, categoryDomain.ErrCategoryNotFound
	}

	summary := toCategorySummary(*category)
	return &summary, nil
}

func (s *CategoryService) Create(ctx context.Context, input dto.CreateCategoryInput) (*dto.CategorySummary, error) {
	name := strings.TrimSpace(input.Name)
	slug := generateSlug(input.Slug, input.Name)

	if err := s.ensureUnique(ctx, uuid.Nil, name, slug); err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	category := &entity.Category{
		ID:          uuid.New(),
		Name:        name,
		Slug:        slug,
		Description: sanitizeDescription(input.Description),
		IsActive:    input.IsActive,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := s.categories.Create(ctx, category); err != nil {
		return nil, sharedErrors.Internal("failed to create category", err)
	}

	summary := toCategorySummary(*category)
	return &summary, nil
}

func (s *CategoryService) Update(ctx context.Context, input dto.UpdateCategoryInput) (*dto.CategorySummary, error) {
	parsedID, err := parseUUID(input.ID, "id")
	if err != nil {
		return nil, err
	}

	category, err := s.categories.FindByID(ctx, parsedID)
	if err != nil {
		return nil, sharedErrors.Internal("failed to load category", err)
	}
	if category == nil {
		return nil, categoryDomain.ErrCategoryNotFound
	}

	name := strings.TrimSpace(input.Name)
	slug := generateSlug(input.Slug, input.Name)
	if err := s.ensureUnique(ctx, category.ID, name, slug); err != nil {
		return nil, err
	}

	category.Name = name
	category.Slug = slug
	category.Description = sanitizeDescription(input.Description)
	category.IsActive = input.IsActive
	category.UpdatedAt = time.Now().UTC()

	if err := s.categories.Update(ctx, category); err != nil {
		return nil, sharedErrors.Internal("failed to update category", err)
	}

	summary := toCategorySummary(*category)
	return &summary, nil
}

func (s *CategoryService) Delete(ctx context.Context, id string) error {
	parsedID, err := parseUUID(id, "id")
	if err != nil {
		return err
	}

	category, err := s.categories.FindByID(ctx, parsedID)
	if err != nil {
		return sharedErrors.Internal("failed to load category", err)
	}
	if category == nil {
		return categoryDomain.ErrCategoryNotFound
	}

	if err := s.categories.Delete(ctx, parsedID); err != nil {
		return sharedErrors.Internal("failed to delete category", err)
	}

	return nil
}

func (s *CategoryService) ensureUnique(ctx context.Context, currentID uuid.UUID, name, slug string) error {
	byName, err := s.categories.FindByName(ctx, name)
	if err != nil {
		return sharedErrors.Internal("failed to check category name", err)
	}
	if byName != nil && byName.ID != currentID {
		return sharedErrors.Conflict("category name already exists")
	}

	bySlug, err := s.categories.FindBySlug(ctx, slug)
	if err != nil {
		return sharedErrors.Internal("failed to check category slug", err)
	}
	if bySlug != nil && bySlug.ID != currentID {
		return sharedErrors.Conflict("category slug already exists")
	}

	return nil
}

func toCategorySummary(category entity.Category) dto.CategorySummary {
	return dto.CategorySummary{
		ID:          category.ID.String(),
		Name:        category.Name,
		Slug:        category.Slug,
		Description: category.Description,
		IsActive:    category.IsActive,
		CreatedAt:   category.CreatedAt,
		UpdatedAt:   category.UpdatedAt,
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
