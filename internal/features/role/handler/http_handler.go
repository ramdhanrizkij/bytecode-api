package handler

import (
	"errors"

	"github.com/gofiber/fiber/v3"
	"go.uber.org/zap"

	"github.com/ramdhanrizkij/bytecode-api/internal/core/middleware"
	"github.com/ramdhanrizkij/bytecode-api/internal/features/role/domain"
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

// GetAll godoc
// @Summary List roles
// @Description Returns paginated roles with their permissions.
// @Tags Roles
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param per_page query int false "Items per page" default(10)
// @Param sort query string false "Sort field" default(created_at)
// @Param order query string false "Sort direction" Enums(asc,desc) default(desc)
// @Param search query string false "Search keyword"
// @Success 200 {object} swaggerdocs.RoleListResponse
// @Failure 401 {object} swaggerdocs.ErrorResponse
// @Failure 403 {object} swaggerdocs.ErrorResponse
// @Failure 500 {object} swaggerdocs.ErrorResponse
// @Router /roles/ [get]
func (h *RoleHTTPHandler) GetAll(c fiber.Ctx) error {
	pq := pagination.NewPaginationQuery(c)
	roles, meta, err := h.service.GetAll(c.Context(), pq)
	if err != nil {
		return h.handleError(c, err)
	}
	return response.SuccessWithPagination(c, "Roles retrieved successfully", roles, meta)
}

// GetByID godoc
// @Summary Get role by ID
// @Description Returns one role by UUID.
// @Tags Roles
// @Produce json
// @Security BearerAuth
// @Param id path string true "Role UUID"
// @Success 200 {object} swaggerdocs.RoleResponse
// @Failure 401 {object} swaggerdocs.ErrorResponse
// @Failure 403 {object} swaggerdocs.ErrorResponse
// @Failure 404 {object} swaggerdocs.ErrorResponse
// @Failure 500 {object} swaggerdocs.ErrorResponse
// @Router /roles/{id} [get]
func (h *RoleHTTPHandler) GetByID(c fiber.Ctx) error {
	id := c.Params("id")
	role, err := h.service.GetByID(c.Context(), id)
	if err != nil {
		return h.handleError(c, err)
	}
	return response.Success(c, "Role retrieved successfully", role)
}

// Create godoc
// @Summary Create role
// @Description Creates a new role.
// @Tags Roles
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param payload body domain.CreateRoleRequest true "Role payload"
// @Success 201 {object} swaggerdocs.RoleResponse
// @Failure 400 {object} swaggerdocs.ErrorResponse
// @Failure 401 {object} swaggerdocs.ErrorResponse
// @Failure 403 {object} swaggerdocs.ErrorResponse
// @Failure 409 {object} swaggerdocs.ErrorResponse
// @Failure 422 {object} swaggerdocs.ValidationErrorResponse
// @Failure 500 {object} swaggerdocs.ErrorResponse
// @Router /roles/ [post]
func (h *RoleHTTPHandler) Create(c fiber.Ctx) error {
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

// Update godoc
// @Summary Update role
// @Description Updates an existing role by UUID.
// @Tags Roles
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Role UUID"
// @Param payload body domain.UpdateRoleRequest true "Role payload"
// @Success 200 {object} swaggerdocs.RoleResponse
// @Failure 400 {object} swaggerdocs.ErrorResponse
// @Failure 401 {object} swaggerdocs.ErrorResponse
// @Failure 403 {object} swaggerdocs.ErrorResponse
// @Failure 404 {object} swaggerdocs.ErrorResponse
// @Failure 422 {object} swaggerdocs.ValidationErrorResponse
// @Failure 500 {object} swaggerdocs.ErrorResponse
// @Router /roles/{id} [put]
func (h *RoleHTTPHandler) Update(c fiber.Ctx) error {
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

// Delete godoc
// @Summary Delete role
// @Description Deletes a role by UUID.
// @Tags Roles
// @Produce json
// @Security BearerAuth
// @Param id path string true "Role UUID"
// @Success 200 {object} swaggerdocs.ErrorResponse
// @Failure 401 {object} swaggerdocs.ErrorResponse
// @Failure 403 {object} swaggerdocs.ErrorResponse
// @Failure 404 {object} swaggerdocs.ErrorResponse
// @Failure 500 {object} swaggerdocs.ErrorResponse
// @Router /roles/{id} [delete]
func (h *RoleHTTPHandler) Delete(c fiber.Ctx) error {
	id := c.Params("id")
	if err := h.service.Delete(c.Context(), id); err != nil {
		return h.handleError(c, err)
	}
	return response.Success(c, "Role deleted successfully", nil)
}

// AssignPermissions godoc
// @Summary Assign permissions to role
// @Description Adds permissions to a role by UUID.
// @Tags Roles
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Role UUID"
// @Param payload body domain.AssignPermissionsRequest true "Permission IDs"
// @Success 200 {object} swaggerdocs.ErrorResponse
// @Failure 400 {object} swaggerdocs.ErrorResponse
// @Failure 401 {object} swaggerdocs.ErrorResponse
// @Failure 403 {object} swaggerdocs.ErrorResponse
// @Failure 404 {object} swaggerdocs.ErrorResponse
// @Failure 422 {object} swaggerdocs.ValidationErrorResponse
// @Failure 500 {object} swaggerdocs.ErrorResponse
// @Router /roles/{id}/permissions [post]
func (h *RoleHTTPHandler) AssignPermissions(c fiber.Ctx) error {
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

// RemovePermissions godoc
// @Summary Remove permissions from role
// @Description Removes permissions from a role by UUID.
// @Tags Roles
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Role UUID"
// @Param payload body domain.AssignPermissionsRequest true "Permission IDs"
// @Success 200 {object} swaggerdocs.ErrorResponse
// @Failure 400 {object} swaggerdocs.ErrorResponse
// @Failure 401 {object} swaggerdocs.ErrorResponse
// @Failure 403 {object} swaggerdocs.ErrorResponse
// @Failure 404 {object} swaggerdocs.ErrorResponse
// @Failure 422 {object} swaggerdocs.ValidationErrorResponse
// @Failure 500 {object} swaggerdocs.ErrorResponse
// @Router /roles/{id}/permissions [delete]
func (h *RoleHTTPHandler) RemovePermissions(c fiber.Ctx) error {
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
func (h *RoleHTTPHandler) handleError(c fiber.Ctx, err error) error {
	var appErr *apperrors.AppError
	if errors.As(err, &appErr) {
		return response.Error(c, appErr.Code, appErr.Message)
	}
	h.log.Error("unexpected error in role handler", zap.Error(err))
	return response.Error(c, fiber.StatusInternalServerError, "internal server error")
}

// ensure middleware import is used (GetCurrentUser available for future RBAC use)
var _ = middleware.GetCurrentUser
