package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/ramdhanrizki/bytecode-api/internal/identity/application/dto"
	identityService "github.com/ramdhanrizki/bytecode-api/internal/identity/application/service"
	httpRequest "github.com/ramdhanrizki/bytecode-api/internal/identity/delivery/http/request"
	httpResponse "github.com/ramdhanrizki/bytecode-api/internal/identity/delivery/http/response"
	sharedAuth "github.com/ramdhanrizki/bytecode-api/internal/shared/auth"
	sharedResponse "github.com/ramdhanrizki/bytecode-api/internal/shared/response"
)

type ProfileHandler struct {
	profileService *identityService.ProfileService
}

func NewProfileHandler(profileService *identityService.ProfileService) *ProfileHandler {
	return &ProfileHandler{profileService: profileService}
}

func (h *ProfileHandler) GetCurrent(c *gin.Context) {
	authUser, ok := sharedAuth.GetAuthenticatedUser(c)
	if !ok {
		sharedResponse.Error(c, validationError(nil))
		return
	}

	output, err := h.profileService.GetCurrent(c.Request.Context(), authUser.ID.String())
	if err != nil {
		sharedResponse.Error(c, err)
		return
	}

	sharedResponse.Success(c, http.StatusOK, "profile fetched successfully", httpResponse.FromUserSummary(*output))
}

func (h *ProfileHandler) UpdateCurrent(c *gin.Context) {
	authUser, ok := sharedAuth.GetAuthenticatedUser(c)
	if !ok {
		sharedResponse.Error(c, validationError(nil))
		return
	}

	var req httpRequest.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		sharedResponse.Error(c, validationError(err))
		return
	}

	output, err := h.profileService.UpdateCurrent(c.Request.Context(), dto.UpdateProfileInput{
		UserID:   authUser.ID.String(),
		FullName: req.FullName,
	})
	if err != nil {
		sharedResponse.Error(c, err)
		return
	}

	sharedResponse.Success(c, http.StatusOK, "profile updated successfully", httpResponse.FromUserSummary(*output))
}
