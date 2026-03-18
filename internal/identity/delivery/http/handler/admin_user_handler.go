package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/ramdhanrizki/bytecode-api/internal/identity/application/dto"
	identityService "github.com/ramdhanrizki/bytecode-api/internal/identity/application/service"
	httpRequest "github.com/ramdhanrizki/bytecode-api/internal/identity/delivery/http/request"
	httpResponse "github.com/ramdhanrizki/bytecode-api/internal/identity/delivery/http/response"
	sharedKernel "github.com/ramdhanrizki/bytecode-api/internal/shared/kernel"
	sharedResponse "github.com/ramdhanrizki/bytecode-api/internal/shared/response"
)

type AdminUserHandler struct {
	service *identityService.AdminUserService
}

func NewAdminUserHandler(service *identityService.AdminUserService) *AdminUserHandler {
	return &AdminUserHandler{service: service}
}

func (h *AdminUserHandler) List(c *gin.Context) {
	query := sharedKernel.PaginationFromRequest(c)
	output, err := h.service.List(c.Request.Context(), dto.ListInput{
		Page:   query.Page,
		Limit:  query.Limit,
		Search: query.Search,
		Sort:   query.Sort,
		Order:  query.Order,
	})
	if err != nil {
		sharedResponse.Error(c, err)
		return
	}

	sharedResponse.Paginated(c, "users fetched successfully", httpResponse.FromUserSummaries(output.Users), toResponseMeta(output.Meta))
}

func (h *AdminUserHandler) Get(c *gin.Context) {
	output, err := h.service.Get(c.Request.Context(), c.Param("id"))
	if err != nil {
		sharedResponse.Error(c, err)
		return
	}

	sharedResponse.Success(c, http.StatusOK, "user fetched successfully", httpResponse.FromUserSummary(*output))
}

func (h *AdminUserHandler) Create(c *gin.Context) {
	var req httpRequest.CreateAdminUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		sharedResponse.Error(c, validationError(err))
		return
	}

	output, err := h.service.Create(c.Request.Context(), dto.CreateAdminUserInput{
		FullName:        req.FullName,
		Email:           req.Email,
		Password:        req.Password,
		IsEmailVerified: req.IsEmailVerified,
		IsActive:        req.IsActive,
		RoleIDs:         req.RoleIDs,
	})
	if err != nil {
		sharedResponse.Error(c, err)
		return
	}

	sharedResponse.Success(c, http.StatusCreated, "user created successfully", httpResponse.FromUserSummary(*output))
}

func (h *AdminUserHandler) Update(c *gin.Context) {
	var req httpRequest.UpdateAdminUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		sharedResponse.Error(c, validationError(err))
		return
	}

	output, err := h.service.Update(c.Request.Context(), dto.UpdateAdminUserInput{
		ID:              c.Param("id"),
		FullName:        req.FullName,
		Email:           req.Email,
		IsEmailVerified: req.IsEmailVerified,
		IsActive:        req.IsActive,
	})
	if err != nil {
		sharedResponse.Error(c, err)
		return
	}

	sharedResponse.Success(c, http.StatusOK, "user updated successfully", httpResponse.FromUserSummary(*output))
}

func (h *AdminUserHandler) Delete(c *gin.Context) {
	if err := h.service.Delete(c.Request.Context(), c.Param("id")); err != nil {
		sharedResponse.Error(c, err)
		return
	}

	sharedResponse.Success(c, http.StatusOK, "user deleted successfully", gin.H{"deleted": true})
}

func (h *AdminUserHandler) AssignRoles(c *gin.Context) {
	var req httpRequest.AssignUserRolesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		sharedResponse.Error(c, validationError(err))
		return
	}

	output, err := h.service.AssignRoles(c.Request.Context(), dto.AssignUserRolesInput{
		UserID:  c.Param("id"),
		RoleIDs: req.RoleIDs,
	})
	if err != nil {
		sharedResponse.Error(c, err)
		return
	}

	sharedResponse.Success(c, http.StatusOK, "user roles updated successfully", httpResponse.FromUserSummary(*output))
}

func toResponseMeta(meta dto.PaginationMeta) sharedResponse.Meta {
	return sharedResponse.Meta{
		Page:       meta.Page,
		Limit:      meta.Limit,
		Total:      meta.Total,
		TotalPages: meta.TotalPages,
	}
}
