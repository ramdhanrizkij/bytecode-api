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

type RoleHandler struct {
	service *identityService.RoleService
}

func NewRoleHandler(service *identityService.RoleService) *RoleHandler {
	return &RoleHandler{service: service}
}

func (h *RoleHandler) List(c *gin.Context) {
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

	sharedResponse.Paginated(c, "roles fetched successfully", httpResponse.FromRoleSummaries(output.Roles), toResponseMeta(output.Meta))
}

func (h *RoleHandler) Get(c *gin.Context) {
	output, err := h.service.Get(c.Request.Context(), c.Param("id"))
	if err != nil {
		sharedResponse.Error(c, err)
		return
	}

	sharedResponse.Success(c, http.StatusOK, "role fetched successfully", httpResponse.FromRoleSummary(*output))
}

func (h *RoleHandler) Create(c *gin.Context) {
	var req httpRequest.CreateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		sharedResponse.Error(c, validationError(err))
		return
	}

	output, err := h.service.Create(c.Request.Context(), dto.CreateRoleInput{
		Name:        req.Name,
		Description: req.Description,
	})
	if err != nil {
		sharedResponse.Error(c, err)
		return
	}

	sharedResponse.Success(c, http.StatusCreated, "role created successfully", httpResponse.FromRoleSummary(*output))
}

func (h *RoleHandler) Update(c *gin.Context) {
	var req httpRequest.UpdateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		sharedResponse.Error(c, validationError(err))
		return
	}

	output, err := h.service.Update(c.Request.Context(), dto.UpdateRoleInput{
		ID:          c.Param("id"),
		Name:        req.Name,
		Description: req.Description,
	})
	if err != nil {
		sharedResponse.Error(c, err)
		return
	}

	sharedResponse.Success(c, http.StatusOK, "role updated successfully", httpResponse.FromRoleSummary(*output))
}

func (h *RoleHandler) Delete(c *gin.Context) {
	if err := h.service.Delete(c.Request.Context(), c.Param("id")); err != nil {
		sharedResponse.Error(c, err)
		return
	}

	sharedResponse.Success(c, http.StatusOK, "role deleted successfully", gin.H{"deleted": true})
}

func (h *RoleHandler) AssignPermissions(c *gin.Context) {
	var req httpRequest.AssignRolePermissionsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		sharedResponse.Error(c, validationError(err))
		return
	}

	output, err := h.service.AssignPermissions(c.Request.Context(), dto.AssignRolePermissionsInput{
		RoleID:        c.Param("id"),
		PermissionIDs: req.PermissionIDs,
	})
	if err != nil {
		sharedResponse.Error(c, err)
		return
	}

	sharedResponse.Success(c, http.StatusOK, "role permissions updated successfully", httpResponse.FromRoleSummary(*output))
}
