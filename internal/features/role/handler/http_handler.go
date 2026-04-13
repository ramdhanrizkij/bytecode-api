package handler

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"

	"github.com/ramdhanrizkij/bytecode-api/internal/features/role/domain"
	"github.com/ramdhanrizkij/bytecode-api/internal/core/middleware"
	apperrors "github.com/ramdhanrizkij/bytecode-api/internal/shared/errors"
	"github.com/ramdhanrizkij/bytecode-api/internal/shared/pagination"
	"github.com/ramdhanrizkij/bytecode-api/internal/shared/response"
	"github.com/ramdhanrizkij/bytecode-api/internal/shared/validator"
)

// RoleHTTPHandler handles HTTP requests for the role feature.
type RoleHTTPHandler struct {
	service domain.RoleService
	log     *zap.Logger
}

// NewRoleHTTPHandler creates a new RoleHTTPHandler.
func NewRoleHTTPHandler(service domain.RoleService, log *zap.Logger) *RoleHTTPHandler {
	return &RoleHTTPHandler{service: service, log: log}
}

// GetAll handles GET /roles
func (h *RoleHTTPHandler) GetAll(c *fiber.Ctx) error {
	pq := pagination.NewPaginationQuery(c)
	roles, meta, err := h.service.GetAll(c.Context(), pq)
	if err != nil {
		return h.handleError(c, err)
	}
	return response.SuccessWithPagination(c, "Roles retrieved successfully", roles, meta)
}

// GetByID handles GET /roles/:id
func (h *RoleHTTPHandler) GetByID(c *fiber.Ctx) error {
	id := c.Params("id")
	role, err := h.service.GetByID(c.Context(), id)
	if err != nil {
		return h.handleError(c, err)
	}
	return response.Success(c, "Role retrieved successfully", role)
}

// Create handles POST /roles
func (h *RoleHTTPHandler) Create(c *fiber.Ctx) error {
	req, err := validator.ParseAndValidate[domain.CreateRoleRequest](c)
	if err != nil {
		return nil
	}
	role, err := h.service.Create(c.Context(), req)
	if err != nil {
		return h.handleError(c, err)
	}
	return response.Created(c, "Role created successfully", role)
}

// Update handles PUT /roles/:id
func (h *RoleHTTPHandler) Update(c *fiber.Ctx) error {
	id := c.Params("id")
	req, err := validator.ParseAndValidate[domain.UpdateRoleRequest](c)
	if err != nil {
		return nil
	}
	role, err := h.service.Update(c.Context(), id, req)
	if err != nil {
		return h.handleError(c, err)
	}
	return response.Success(c, "Role updated successfully", role)
}

// Delete handles DELETE /roles/:id
func (h *RoleHTTPHandler) Delete(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := h.service.Delete(c.Context(), id); err != nil {
		return h.handleError(c, err)
	}
	return response.Success(c, "Role deleted successfully", nil)
}

// AssignPermissions handles POST /roles/:id/permissions
func (h *RoleHTTPHandler) AssignPermissions(c *fiber.Ctx) error {
	id := c.Params("id")
	req, err := validator.ParseAndValidate[domain.AssignPermissionsRequest](c)
	if err != nil {
		return nil
	}
	if err := h.service.AssignPermissions(c.Context(), id, req); err != nil {
		return h.handleError(c, err)
	}
	return response.Success(c, "Permissions assigned to role successfully", nil)
}

// RemovePermissions handles DELETE /roles/:id/permissions
func (h *RoleHTTPHandler) RemovePermissions(c *fiber.Ctx) error {
	id := c.Params("id")
	req, err := validator.ParseAndValidate[domain.AssignPermissionsRequest](c)
	if err != nil {
		return nil
	}
	if err := h.service.RemovePermissions(c.Context(), id, req); err != nil {
		return h.handleError(c, err)
	}
	return response.Success(c, "Permissions removed from role successfully", nil)
}

// handleError maps AppError to HTTP response codes.
func (h *RoleHTTPHandler) handleError(c *fiber.Ctx, err error) error {
	var appErr *apperrors.AppError
	if errors.As(err, &appErr) {
		return response.Error(c, appErr.Code, appErr.Message)
	}
	h.log.Error("unexpected error in role handler", zap.Error(err))
	return response.Error(c, fiber.StatusInternalServerError, "internal server error")
}

// ensure middleware import is used (GetCurrentUser available for future RBAC use)
var _ = middleware.GetCurrentUser
