package handler

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/ramdhanrizki/bytecode-api/internal/product/application/dto"
	productService "github.com/ramdhanrizki/bytecode-api/internal/product/application/service"
	httpRequest "github.com/ramdhanrizki/bytecode-api/internal/product/delivery/http/request"
	httpResponse "github.com/ramdhanrizki/bytecode-api/internal/product/delivery/http/response"
	sharedKernel "github.com/ramdhanrizki/bytecode-api/internal/shared/kernel"
	sharedResponse "github.com/ramdhanrizki/bytecode-api/internal/shared/response"
)

type ProductHandler struct {
	service *productService.ProductService
}

func NewProductHandler(service *productService.ProductService) *ProductHandler {
	return &ProductHandler{service: service}
}

func (h *ProductHandler) List(c *gin.Context) {
	query := sharedKernel.PaginationFromRequest(c)
	output, err := h.service.List(c.Request.Context(), dto.ListInput{
		Page:       query.Page,
		Limit:      query.Limit,
		Search:     query.Search,
		Sort:       query.Sort,
		Order:      query.Order,
		CategoryID: stringPointer(strings.TrimSpace(c.Query("category_id"))),
		IsActive:   parseBoolQuery(c.Query("is_active")),
	})
	if err != nil {
		sharedResponse.Error(c, err)
		return
	}

	sharedResponse.Paginated(c, "products fetched successfully", httpResponse.FromProductSummaries(output.Products), sharedResponse.Meta{
		Page:       output.Meta.Page,
		Limit:      output.Meta.Limit,
		Total:      output.Meta.Total,
		TotalPages: output.Meta.TotalPages,
	})
}

func (h *ProductHandler) Get(c *gin.Context) {
	output, err := h.service.Get(c.Request.Context(), c.Param("id"))
	if err != nil {
		sharedResponse.Error(c, err)
		return
	}

	sharedResponse.Success(c, http.StatusOK, "product fetched successfully", httpResponse.FromProductSummary(*output))
}

func (h *ProductHandler) Create(c *gin.Context) {
	var req httpRequest.UpsertProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		sharedResponse.Error(c, validationError(err))
		return
	}

	output, err := h.service.Create(c.Request.Context(), dto.CreateProductInput{
		CategoryID:  req.CategoryID,
		Name:        req.Name,
		Slug:        req.Slug,
		Description: req.Description,
		SKU:         req.SKU,
		Price:       req.Price,
		Stock:       req.Stock,
		IsActive:    req.IsActive,
	})
	if err != nil {
		sharedResponse.Error(c, err)
		return
	}

	sharedResponse.Success(c, http.StatusCreated, "product created successfully", httpResponse.FromProductSummary(*output))
}

func (h *ProductHandler) Update(c *gin.Context) {
	var req httpRequest.UpsertProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		sharedResponse.Error(c, validationError(err))
		return
	}

	output, err := h.service.Update(c.Request.Context(), dto.UpdateProductInput{
		ID:          c.Param("id"),
		CategoryID:  req.CategoryID,
		Name:        req.Name,
		Slug:        req.Slug,
		Description: req.Description,
		SKU:         req.SKU,
		Price:       req.Price,
		Stock:       req.Stock,
		IsActive:    req.IsActive,
	})
	if err != nil {
		sharedResponse.Error(c, err)
		return
	}

	sharedResponse.Success(c, http.StatusOK, "product updated successfully", httpResponse.FromProductSummary(*output))
}

func (h *ProductHandler) Delete(c *gin.Context) {
	if err := h.service.Delete(c.Request.Context(), c.Param("id")); err != nil {
		sharedResponse.Error(c, err)
		return
	}

	sharedResponse.Success(c, http.StatusOK, "product deleted successfully", gin.H{"deleted": true})
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

func stringPointer(value string) *string {
	if value == "" {
		return nil
	}
	return &value
}
