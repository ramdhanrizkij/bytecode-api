package handler

import (
	"errors"

	"github.com/gofiber/fiber/v3"
	"go.uber.org/zap"

	"github.com/ramdhanrizkij/bytecode-api/internal/features/auth/domain"
	apperrors "github.com/ramdhanrizkij/bytecode-api/internal/shared/errors"
	"github.com/ramdhanrizkij/bytecode-api/internal/shared/response"
	"github.com/ramdhanrizkij/bytecode-api/internal/shared/validator"
)

// AuthHTTPHandler handles HTTP requests for the auth feature.
// It is the only layer allowed to import Fiber.
type AuthHTTPHandler struct {
	service domain.AuthService
	log     *zap.Logger
}

// NewAuthHTTPHandler creates a new AuthHTTPHandler.
func NewAuthHTTPHandler(service domain.AuthService, log *zap.Logger) *AuthHTTPHandler {
	return &AuthHTTPHandler{service: service, log: log}
}

// Register godoc
// @Summary Register a new user
// @Description Creates a user account and returns access and refresh tokens.
// @Tags Auth
// @Accept json
// @Produce json
// @Param payload body domain.RegisterRequest true "Registration payload"
// @Success 201 {object} swaggerdocs.AuthResponse
// @Failure 400 {object} swaggerdocs.ErrorResponse
// @Failure 409 {object} swaggerdocs.ErrorResponse
// @Failure 422 {object} swaggerdocs.ValidationErrorResponse
// @Failure 500 {object} swaggerdocs.ErrorResponse
// @Router /auth/register [post]
func (h *AuthHTTPHandler) Register(c fiber.Ctx) error {
	req, err := validator.ParseAndValidate[domain.RegisterRequest](c)
	if err != nil {
		// ParseAndValidate already sent the 422 response.
		return nil
	}

	resp, err := h.service.Register(c.Context(), req)
	if err != nil {
		return h.handleError(c, err)
	}

	return response.Created(c, "User registered successfully", resp)
}

// Login godoc
// @Summary Login
// @Description Authenticates a user and returns access and refresh tokens.
// @Tags Auth
// @Accept json
// @Produce json
// @Param payload body domain.LoginRequest true "Login payload"
// @Success 200 {object} swaggerdocs.AuthResponse
// @Failure 400 {object} swaggerdocs.ErrorResponse
// @Failure 401 {object} swaggerdocs.ErrorResponse
// @Failure 422 {object} swaggerdocs.ValidationErrorResponse
// @Failure 500 {object} swaggerdocs.ErrorResponse
// @Router /auth/login [post]
func (h *AuthHTTPHandler) Login(c fiber.Ctx) error {
	req, err := validator.ParseAndValidate[domain.LoginRequest](c)
	if err != nil {
		return nil
	}

	resp, err := h.service.Login(c.Context(), req)
	if err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, "Login successful", resp)
}

// Refresh godoc
// @Summary Refresh access token
// @Description Exchanges a valid refresh token for a new token pair.
// @Tags Auth
// @Accept json
// @Produce json
// @Param payload body domain.RefreshTokenRequest true "Refresh token payload"
// @Success 200 {object} swaggerdocs.TokenResponse
// @Failure 400 {object} swaggerdocs.ErrorResponse
// @Failure 401 {object} swaggerdocs.ErrorResponse
// @Failure 422 {object} swaggerdocs.ValidationErrorResponse
// @Failure 500 {object} swaggerdocs.ErrorResponse
// @Router /auth/refresh [post]
func (h *AuthHTTPHandler) Refresh(c fiber.Ctx) error {
	req, err := validator.ParseAndValidate[domain.RefreshTokenRequest](c)
	if err != nil {
		return nil
	}

	resp, err := h.service.Refresh(c.Context(), req)
	if err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, "Token refreshed successfully", resp)
}

// Logout godoc
// @Summary Logout
// @Description Revokes a refresh token.
// @Tags Auth
// @Accept json
// @Produce json
// @Param payload body domain.LogoutRequest true "Logout payload"
// @Success 200 {object} swaggerdocs.ErrorResponse
// @Failure 400 {object} swaggerdocs.ErrorResponse
// @Failure 422 {object} swaggerdocs.ValidationErrorResponse
// @Failure 500 {object} swaggerdocs.ErrorResponse
// @Router /auth/logout [post]
func (h *AuthHTTPHandler) Logout(c fiber.Ctx) error {
	req, err := validator.ParseAndValidate[domain.LogoutRequest](c)
	if err != nil {
		return nil
	}

	if err := h.service.Logout(c.Context(), req); err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, "Logout successful", nil)
}

// handleError maps AppError values to appropriate HTTP responses.
func (h *AuthHTTPHandler) handleError(c fiber.Ctx, err error) error {
	var appErr *apperrors.AppError
	if errors.As(err, &appErr) {
		return response.Error(c, appErr.Code, appErr.Message)
	}

	h.log.Error("unexpected error in auth handler", zap.Error(err))
	return response.Error(c, fiber.StatusInternalServerError, "internal server error")
}
