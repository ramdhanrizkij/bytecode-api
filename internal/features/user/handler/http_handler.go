package handler

import (
	"errors"

	"github.com/gofiber/fiber/v3"
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

// GetAll godoc
// @Summary List users
// @Description Returns paginated users.
// @Tags Users
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param per_page query int false "Items per page" default(10)
// @Param sort query string false "Sort field" default(created_at)
// @Param order query string false "Sort direction" Enums(asc,desc) default(desc)
// @Param search query string false "Search keyword"
// @Success 200 {object} swaggerdocs.UserListResponse
// @Failure 401 {object} swaggerdocs.ErrorResponse
// @Failure 403 {object} swaggerdocs.ErrorResponse
// @Failure 500 {object} swaggerdocs.ErrorResponse
// @Router /users/ [get]
func (h *UserHTTPHandler) GetAll(c fiber.Ctx) error {
	pq := pagination.NewPaginationQuery(c)
	users, meta, err := h.service.GetAll(c.Context(), pq)
	if err != nil {
		return h.handleError(c, err)
	}
	return response.SuccessWithPagination(c, "Users retrieved successfully", users, meta)
}

// GetByID godoc
// @Summary Get user by ID
// @Description Returns one user by UUID.
// @Tags Users
// @Produce json
// @Security BearerAuth
// @Param id path string true "User UUID"
// @Success 200 {object} swaggerdocs.UserResponse
// @Failure 401 {object} swaggerdocs.ErrorResponse
// @Failure 403 {object} swaggerdocs.ErrorResponse
// @Failure 404 {object} swaggerdocs.ErrorResponse
// @Failure 500 {object} swaggerdocs.ErrorResponse
// @Router /users/{id} [get]
func (h *UserHTTPHandler) GetByID(c fiber.Ctx) error {
	id := c.Params("id")
	user, err := h.service.GetByID(c.Context(), id)
	if err != nil {
		return h.handleError(c, err)
	}
	return response.Success(c, "User retrieved successfully", user)
}

// GetMe godoc
// @Summary Get current user
// @Description Returns the profile for the authenticated user.
// @Tags Users
// @Produce json
// @Security BearerAuth
// @Success 200 {object} swaggerdocs.UserResponse
// @Failure 401 {object} swaggerdocs.ErrorResponse
// @Failure 404 {object} swaggerdocs.ErrorResponse
// @Failure 500 {object} swaggerdocs.ErrorResponse
// @Router /users/me [get]
func (h *UserHTTPHandler) GetMe(c fiber.Ctx) error {
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

// UpdateProfile godoc
// @Summary Update current user profile
// @Description Updates the authenticated user's profile.
// @Tags Users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param payload body domain.UpdateProfileRequest true "Profile payload"
// @Success 200 {object} swaggerdocs.UserResponse
// @Failure 400 {object} swaggerdocs.ErrorResponse
// @Failure 401 {object} swaggerdocs.ErrorResponse
// @Failure 422 {object} swaggerdocs.ValidationErrorResponse
// @Failure 500 {object} swaggerdocs.ErrorResponse
// @Router /users/me/profile [put]
func (h *UserHTTPHandler) UpdateProfile(c fiber.Ctx) error {
	claims := middleware.GetCurrentUser(c)
	if claims == nil {
		return response.Error(c, fiber.StatusUnauthorized, "unauthorized")
	}

	req, err := validator.ParseAndValidate[domain.UpdateProfileRequest](c)
	if err != nil {
		return nil
	}

	user, err := h.service.UpdateProfile(c.Context(), claims.UserID, req)
	if err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, "Profile updated successfully", user)
}

// GetMePermissions godoc
// @Summary Get current user permissions
// @Description Returns permission names assigned to the authenticated user's role.
// @Tags Users
// @Produce json
// @Security BearerAuth
// @Success 200 {object} swaggerdocs.UserPermissionListResponse
// @Failure 401 {object} swaggerdocs.ErrorResponse
// @Failure 500 {object} swaggerdocs.ErrorResponse
// @Router /users/me/permissions [get]
func (h *UserHTTPHandler) GetMePermissions(c fiber.Ctx) error {
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

// Create godoc
// @Summary Create user
// @Description Creates a new user.
// @Tags Users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param payload body domain.CreateUserRequest true "User payload"
// @Success 201 {object} swaggerdocs.UserResponse
// @Failure 400 {object} swaggerdocs.ErrorResponse
// @Failure 401 {object} swaggerdocs.ErrorResponse
// @Failure 403 {object} swaggerdocs.ErrorResponse
// @Failure 409 {object} swaggerdocs.ErrorResponse
// @Failure 422 {object} swaggerdocs.ValidationErrorResponse
// @Failure 500 {object} swaggerdocs.ErrorResponse
// @Router /users/ [post]
func (h *UserHTTPHandler) Create(c fiber.Ctx) error {
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

// Update godoc
// @Summary Update user
// @Description Updates an existing user by UUID.
// @Tags Users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "User UUID"
// @Param payload body domain.UpdateUserRequest true "User payload"
// @Success 200 {object} swaggerdocs.UserResponse
// @Failure 400 {object} swaggerdocs.ErrorResponse
// @Failure 401 {object} swaggerdocs.ErrorResponse
// @Failure 403 {object} swaggerdocs.ErrorResponse
// @Failure 404 {object} swaggerdocs.ErrorResponse
// @Failure 422 {object} swaggerdocs.ValidationErrorResponse
// @Failure 500 {object} swaggerdocs.ErrorResponse
// @Router /users/{id} [put]
func (h *UserHTTPHandler) Update(c fiber.Ctx) error {
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

// Delete godoc
// @Summary Delete user
// @Description Deletes a user by UUID.
// @Tags Users
// @Produce json
// @Security BearerAuth
// @Param id path string true "User UUID"
// @Success 200 {object} swaggerdocs.ErrorResponse
// @Failure 401 {object} swaggerdocs.ErrorResponse
// @Failure 403 {object} swaggerdocs.ErrorResponse
// @Failure 404 {object} swaggerdocs.ErrorResponse
// @Failure 500 {object} swaggerdocs.ErrorResponse
// @Router /users/{id} [delete]
func (h *UserHTTPHandler) Delete(c fiber.Ctx) error {
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
func (h *UserHTTPHandler) handleError(c fiber.Ctx, err error) error {
	var appErr *apperrors.AppError
	if errors.As(err, &appErr) {
		return response.Error(c, appErr.Code, appErr.Message)
	}

	h.log.Error("unexpected user handler error", zap.Error(err))
	return response.Error(c, fiber.StatusInternalServerError, "internal server error")
}
