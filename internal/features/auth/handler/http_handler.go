package handler

import (
	"errors"

	"github.com/gofiber/fiber/v2"
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

// Register handles POST /auth/register.
func (h *AuthHTTPHandler) Register(c *fiber.Ctx) error {
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

// Login handles POST /auth/login.
func (h *AuthHTTPHandler) Login(c *fiber.Ctx) error {
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

// handleError maps AppError values to appropriate HTTP responses.
func (h *AuthHTTPHandler) handleError(c *fiber.Ctx, err error) error {
	var appErr *apperrors.AppError
	if errors.As(err, &appErr) {
		return response.Error(c, appErr.Code, appErr.Message)
	}

	h.log.Error("unexpected error in auth handler", zap.Error(err))
	return response.Error(c, fiber.StatusInternalServerError, "internal server error")
}
