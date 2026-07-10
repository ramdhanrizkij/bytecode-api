package handler

import (
	"errors"

	"github.com/gofiber/fiber/v3"
	"go.uber.org/zap"

	"github.com/ramdhanrizkij/bytecode-api/internal/features/permission/domain"
	apperrors "github.com/ramdhanrizkij/bytecode-api/internal/shared/errors"
	"github.com/ramdhanrizkij/bytecode-api/internal/shared/pagination"
	"github.com/ramdhanrizkij/bytecode-api/internal/shared/response"
	"github.com/ramdhanrizkij/bytecode-api/internal/shared/validator"
)

// PermissionHTTPHandler handles HTTP requests for the permission feature.
type PermissionHTTPHandler struct {
	service domain.PermissionService
	log     *zap.Logger
}

// NewPermissionHTTPHandler creates a new PermissionHTTPHandler.
func NewPermissionHTTPHandler(service domain.PermissionService, log *zap.Logger) *PermissionHTTPHandler {
	return &PermissionHTTPHandler{service: service, log: log}
}

// GetAll godoc
// @Summary List permissions
// @Description Returns paginated permissions.
// @Tags Permissions
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param per_page query int false "Items per page" default(10)
// @Param sort query string false "Sort field" default(created_at)
// @Param order query string false "Sort direction" Enums(asc,desc) default(desc)
// @Param search query string false "Search keyword"
// @Success 200 {object} swaggerdocs.PermissionListResponse
// @Failure 401 {object} swaggerdocs.ErrorResponse
// @Failure 403 {object} swaggerdocs.ErrorResponse
// @Failure 500 {object} swaggerdocs.ErrorResponse
// @Router /permissions/ [get]
func (h *PermissionHTTPHandler) GetAll(c fiber.Ctx) error {
	pq := pagination.NewPaginationQuery(c)
	perms, meta, err := h.service.GetAll(c.Context(), pq)
	if err != nil {
		return h.handleError(c, err)
	}
	return response.SuccessWithPagination(c, "Permissions retrieved successfully", perms, meta)
}

// GetByID godoc
// @Summary Get permission by ID
// @Description Returns one permission by UUID.
// @Tags Permissions
// @Produce json
// @Security BearerAuth
// @Param id path string true "Permission UUID"
// @Success 200 {object} swaggerdocs.PermissionResponse
// @Failure 401 {object} swaggerdocs.ErrorResponse
// @Failure 403 {object} swaggerdocs.ErrorResponse
// @Failure 404 {object} swaggerdocs.ErrorResponse
// @Failure 500 {object} swaggerdocs.ErrorResponse
// @Router /permissions/{id} [get]
func (h *PermissionHTTPHandler) GetByID(c fiber.Ctx) error {
	id := c.Params("id")
	perm, err := h.service.GetByID(c.Context(), id)
	if err != nil {
		return h.handleError(c, err)
	}
	return response.Success(c, "Permission retrieved successfully", perm)
}

// Create godoc
// @Summary Create permission
// @Description Creates a new permission.
// @Tags Permissions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param payload body domain.CreatePermissionRequest true "Permission payload"
// @Success 201 {object} swaggerdocs.PermissionResponse
// @Failure 400 {object} swaggerdocs.ErrorResponse
// @Failure 401 {object} swaggerdocs.ErrorResponse
// @Failure 403 {object} swaggerdocs.ErrorResponse
// @Failure 409 {object} swaggerdocs.ErrorResponse
// @Failure 422 {object} swaggerdocs.ValidationErrorResponse
// @Failure 500 {object} swaggerdocs.ErrorResponse
// @Router /permissions/ [post]
func (h *PermissionHTTPHandler) Create(c fiber.Ctx) error {
	req, err := validator.ParseAndValidate[domain.CreatePermissionRequest](c)
	if err != nil {
		return nil
	}
	perm, err := h.service.Create(c.Context(), req)
	if err != nil {
		return h.handleError(c, err)
	}
	return response.Created(c, "Permission created successfully", perm)
}

// Update godoc
// @Summary Update permission
// @Description Updates an existing permission by UUID.
// @Tags Permissions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Permission UUID"
// @Param payload body domain.UpdatePermissionRequest true "Permission payload"
// @Success 200 {object} swaggerdocs.PermissionResponse
// @Failure 400 {object} swaggerdocs.ErrorResponse
// @Failure 401 {object} swaggerdocs.ErrorResponse
// @Failure 403 {object} swaggerdocs.ErrorResponse
// @Failure 404 {object} swaggerdocs.ErrorResponse
// @Failure 422 {object} swaggerdocs.ValidationErrorResponse
// @Failure 500 {object} swaggerdocs.ErrorResponse
// @Router /permissions/{id} [put]
func (h *PermissionHTTPHandler) Update(c fiber.Ctx) error {
	id := c.Params("id")
	req, err := validator.ParseAndValidate[domain.UpdatePermissionRequest](c)
	if err != nil {
		return nil
	}
	perm, err := h.service.Update(c.Context(), id, req)
	if err != nil {
		return h.handleError(c, err)
	}
	return response.Success(c, "Permission updated successfully", perm)
}

// Delete godoc
// @Summary Delete permission
// @Description Deletes a permission by UUID.
// @Tags Permissions
// @Produce json
// @Security BearerAuth
// @Param id path string true "Permission UUID"
// @Success 200 {object} swaggerdocs.ErrorResponse
// @Failure 401 {object} swaggerdocs.ErrorResponse
// @Failure 403 {object} swaggerdocs.ErrorResponse
// @Failure 404 {object} swaggerdocs.ErrorResponse
// @Failure 500 {object} swaggerdocs.ErrorResponse
// @Router /permissions/{id} [delete]
func (h *PermissionHTTPHandler) Delete(c fiber.Ctx) error {
	id := c.Params("id")
	if err := h.service.Delete(c.Context(), id); err != nil {
		return h.handleError(c, err)
	}
	return response.Success(c, "Permission deleted successfully", nil)
}

// handleError maps AppError to HTTP response codes.
func (h *PermissionHTTPHandler) handleError(c fiber.Ctx, err error) error {
	var appErr *apperrors.AppError
	if errors.As(err, &appErr) {
		return response.Error(c, appErr.Code, appErr.Message)
	}
	h.log.Error("unexpected error in permission handler", zap.Error(err))
	return response.Error(c, fiber.StatusInternalServerError, "internal server error")
}
