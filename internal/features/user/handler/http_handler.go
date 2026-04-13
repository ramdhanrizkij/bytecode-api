package handler

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"

	"github.com/ramdhanrizkij/bytecode-api/internal/core/middleware"
	"github.com/ramdhanrizkij/bytecode-api/internal/features/user/domain"
	apperrors "github.com/ramdhanrizkij/bytecode-api/internal/shared/errors"
	"github.com/ramdhanrizkij/bytecode-api/internal/shared/pagination"
	"github.com/ramdhanrizkij/bytecode-api/internal/shared/response"
	"github.com/ramdhanrizkij/bytecode-api/internal/shared/validator"
)

// UserHTTPHandler handles HTTP requests for user management.
type UserHTTPHandler struct {
	service domain.UserService
	log     *zap.Logger
}

// NewUserHTTPHandler creates a new UserHTTPHandler instance.
func NewUserHTTPHandler(service domain.UserService, log *zap.Logger) *UserHTTPHandler {
	return &UserHTTPHandler{service: service, log: log}
}

// GetAll handles GET /users
func (h *UserHTTPHandler) GetAll(c *fiber.Ctx) error {
	pq := pagination.NewPaginationQuery(c)
	users, meta, err := h.service.GetAll(c.Context(), pq)
	if err != nil {
		return h.handleError(c, err)
	}
	return response.SuccessWithPagination(c, "Users retrieved successfully", users, meta)
}

// GetByID handles GET /users/:id
func (h *UserHTTPHandler) GetByID(c *fiber.Ctx) error {
	id := c.Params("id")
	user, err := h.service.GetByID(c.Context(), id)
	if err != nil {
		return h.handleError(c, err)
	}
	return response.Success(c, "User retrieved successfully", user)
}

// GetMe handles GET /users/me
func (h *UserHTTPHandler) GetMe(c *fiber.Ctx) error {
	claims := middleware.GetCurrentUser(c)
	if claims == nil {
		return response.Error(c, fiber.StatusUnauthorized, "unauthorized")
	}

	user, err := h.service.GetByID(c.Context(), claims.UserID)
	if err != nil {
		return h.handleError(c, err)
	}
	return response.Success(c, "Current user profile retrieved successfully", user)
}

// GetMePermissions handles GET /users/me/permissions
func (h *UserHTTPHandler) GetMePermissions(c *fiber.Ctx) error {
	claims := middleware.GetCurrentUser(c)
	if claims == nil {
		return response.Error(c, fiber.StatusUnauthorized, "unauthorized")
	}

	perms, err := h.service.GetPermissions(c.Context(), claims.UserID)
	if err != nil {
		return h.handleError(c, err)
	}
	return response.Success(c, "User permissions retrieved successfully", perms)
}

// Create handles POST /users
func (h *UserHTTPHandler) Create(c *fiber.Ctx) error {
	req, err := validator.ParseAndValidate[domain.CreateUserRequest](c)
	if err != nil {
		return nil // validator already returns the response
	}

	user, err := h.service.Create(c.Context(), req)
	if err != nil {
		return h.handleError(c, err)
	}
	return response.Created(c, "User created successfully", user)
}

// Update handles PUT /users/:id
func (h *UserHTTPHandler) Update(c *fiber.Ctx) error {
	id := c.Params("id")
	req, err := validator.ParseAndValidate[domain.UpdateUserRequest](c)
	if err != nil {
		return nil
	}

	user, err := h.service.Update(c.Context(), id, req)
	if err != nil {
		return h.handleError(c, err)
	}
	return response.Success(c, "User updated successfully", user)
}

// Delete handles DELETE /users/:id
func (h *UserHTTPHandler) Delete(c *fiber.Ctx) error {
	id := c.Params("id")
	claims := middleware.GetCurrentUser(c)
	if claims == nil {
		return response.Error(c, fiber.StatusUnauthorized, "unauthorized")
	}

	if err := h.service.Delete(c.Context(), claims.UserID, id); err != nil {
		return h.handleError(c, err)
	}
	return response.Success(c, "User deleted successfully", nil)
}

// handleError maps application errors to HTTP response codes.
func (h *UserHTTPHandler) handleError(c *fiber.Ctx, err error) error {
	var appErr *apperrors.AppError
	if errors.As(err, &appErr) {
		return response.Error(c, appErr.Code, appErr.Message)
	}

	h.log.Error("unexpected user handler error", zap.Error(err))
	return response.Error(c, fiber.StatusInternalServerError, "internal server error")
}
