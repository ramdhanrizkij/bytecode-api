package handler

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/ramdhanrizki/bytecode-api/internal/category/application/dto"
	categoryService "github.com/ramdhanrizki/bytecode-api/internal/category/application/service"
	httpRequest "github.com/ramdhanrizki/bytecode-api/internal/category/delivery/http/request"
	httpResponse "github.com/ramdhanrizki/bytecode-api/internal/category/delivery/http/response"
	sharedKernel "github.com/ramdhanrizki/bytecode-api/internal/shared/kernel"
	sharedResponse "github.com/ramdhanrizki/bytecode-api/internal/shared/response"
)

type CategoryHandler struct {
	service *categoryService.CategoryService
}

func NewCategoryHandler(service *categoryService.CategoryService) *CategoryHandler {
	return &CategoryHandler{service: service}
}

func (h *CategoryHandler) List(c *gin.Context) {
	query := sharedKernel.PaginationFromRequest(c)
	output, err := h.service.List(c.Request.Context(), dto.ListInput{
		Page:     query.Page,
		Limit:    query.Limit,
		Search:   query.Search,
		Sort:     query.Sort,
		Order:    query.Order,
		IsActive: parseBoolQuery(c.Query("is_active")),
	})
	if err != nil {
		sharedResponse.Error(c, err)
		return
	}

	sharedResponse.Paginated(c, "categories fetched successfully", httpResponse.FromCategorySummaries(output.Categories), sharedResponse.Meta{
		Page:       output.Meta.Page,
		Limit:      output.Meta.Limit,
		Total:      output.Meta.Total,
		TotalPages: output.Meta.TotalPages,
	})
}

func (h *CategoryHandler) Get(c *gin.Context) {
	output, err := h.service.Get(c.Request.Context(), c.Param("id"))
	if err != nil {
		sharedResponse.Error(c, err)
		return
	}

	sharedResponse.Success(c, http.StatusOK, "category fetched successfully", httpResponse.FromCategorySummary(*output))
}

func (h *CategoryHandler) Create(c *gin.Context) {
	var req httpRequest.UpsertCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		sharedResponse.Error(c, validationError(err))
		return
	}

	output, err := h.service.Create(c.Request.Context(), dto.CreateCategoryInput{
		Name:        req.Name,
		Slug:        req.Slug,
		Description: req.Description,
		IsActive:    req.IsActive,
	})
	if err != nil {
		sharedResponse.Error(c, err)
		return
	}

	sharedResponse.Success(c, http.StatusCreated, "category created successfully", httpResponse.FromCategorySummary(*output))
}

func (h *CategoryHandler) Update(c *gin.Context) {
	var req httpRequest.UpsertCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		sharedResponse.Error(c, validationError(err))
		return
	}

	output, err := h.service.Update(c.Request.Context(), dto.UpdateCategoryInput{
		ID:          c.Param("id"),
		Name:        req.Name,
		Slug:        req.Slug,
		Description: req.Description,
		IsActive:    req.IsActive,
	})
	if err != nil {
		sharedResponse.Error(c, err)
		return
	}

	sharedResponse.Success(c, http.StatusOK, "category updated successfully", httpResponse.FromCategorySummary(*output))
}

func (h *CategoryHandler) Delete(c *gin.Context) {
	if err := h.service.Delete(c.Request.Context(), c.Param("id")); err != nil {
		sharedResponse.Error(c, err)
		return
	}

	sharedResponse.Success(c, http.StatusOK, "category deleted successfully", gin.H{"deleted": true})
}

func parseBoolQuery(raw string) *bool {
	value := strings.TrimSpace(strings.ToLower(raw))
	if value == "" {
		return nil
	}
	parsed := value == "true" || value == "1"
	if value == "false" || value == "0" {
		parsed = false
	}
	return &parsed
}
