package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/ramdhanrizki/bytecode-api/internal/identity/application/dto"
	identityService "github.com/ramdhanrizki/bytecode-api/internal/identity/application/service"
	httpRequest "github.com/ramdhanrizki/bytecode-api/internal/identity/delivery/http/request"
	httpResponse "github.com/ramdhanrizki/bytecode-api/internal/identity/delivery/http/response"
	sharedResponse "github.com/ramdhanrizki/bytecode-api/internal/shared/response"
)

type AuthHandler struct {
	authService *identityService.AuthService
}

func NewAuthHandler(authService *identityService.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req httpRequest.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		sharedResponse.Error(c, validationError(err))
		return
	}

	output, err := h.authService.Register(c.Request.Context(), dto.RegisterInput{
		FullName: req.FullName,
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		sharedResponse.Error(c, err)
		return
	}

	sharedResponse.Success(c, http.StatusCreated, "user registered successfully", httpResponse.FromRegisterOutput(*output))
}

func (h *AuthHandler) VerifyEmail(c *gin.Context) {
	var req httpRequest.VerifyEmailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		if token := c.Query("token"); token != "" {
			req.Token = token
		} else {
			sharedResponse.Error(c, validationError(err))
			return
		}
	}

	if err := h.authService.VerifyEmail(c.Request.Context(), dto.VerifyEmailInput{Token: req.Token}); err != nil {
		sharedResponse.Error(c, err)
		return
	}

	sharedResponse.Success(c, http.StatusOK, "email verified successfully", gin.H{"verified": true})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req httpRequest.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		sharedResponse.Error(c, validationError(err))
		return
	}

	output, err := h.authService.Login(c.Request.Context(), dto.LoginInput{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		sharedResponse.Error(c, err)
		return
	}

	sharedResponse.Success(c, http.StatusOK, "login successful", httpResponse.FromAuthOutput(*output))
}

func (h *AuthHandler) Refresh(c *gin.Context) {
	var req httpRequest.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		sharedResponse.Error(c, validationError(err))
		return
	}

	output, err := h.authService.RefreshToken(c.Request.Context(), dto.RefreshTokenInput{
		RefreshToken: req.RefreshToken,
	})
	if err != nil {
		sharedResponse.Error(c, err)
		return
	}

	sharedResponse.Success(c, http.StatusOK, "token refreshed successfully", httpResponse.FromAuthOutput(*output))
}
