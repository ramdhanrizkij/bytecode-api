package handler

import (
	"errors"

	"github.com/gofiber/fiber/v2"
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

// GetAll handles GET /permissions
func (h *PermissionHTTPHandler) GetAll(c *fiber.Ctx) error {
	pq := pagination.NewPaginationQuery(c)
	perms, meta, err := h.service.GetAll(c.Context(), pq)
	if err != nil {
		return h.handleError(c, err)
	}
	return response.SuccessWithPagination(c, "Permissions retrieved successfully", perms, meta)
}

// GetByID handles GET /permissions/:id
func (h *PermissionHTTPHandler) GetByID(c *fiber.Ctx) error {
	id := c.Params("id")
	perm, err := h.service.GetByID(c.Context(), id)
	if err != nil {
		return h.handleError(c, err)
	}
	return response.Success(c, "Permission retrieved successfully", perm)
}

// Create handles POST /permissions
func (h *PermissionHTTPHandler) Create(c *fiber.Ctx) error {
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

// Update handles PUT /permissions/:id
func (h *PermissionHTTPHandler) Update(c *fiber.Ctx) error {
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

// Delete handles DELETE /permissions/:id
func (h *PermissionHTTPHandler) Delete(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := h.service.Delete(c.Context(), id); err != nil {
		return h.handleError(c, err)
	}
	return response.Success(c, "Permission deleted successfully", nil)
}

// handleError maps AppError to HTTP response codes.
func (h *PermissionHTTPHandler) handleError(c *fiber.Ctx, err error) error {
	var appErr *apperrors.AppError
	if errors.As(err, &appErr) {
		return response.Error(c, appErr.Code, appErr.Message)
	}
	h.log.Error("unexpected error in permission handler", zap.Error(err))
	return response.Error(c, fiber.StatusInternalServerError, "internal server error")
}
