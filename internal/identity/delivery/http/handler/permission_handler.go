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

type PermissionHandler struct {
	service *identityService.PermissionService
}

func NewPermissionHandler(service *identityService.PermissionService) *PermissionHandler {
	return &PermissionHandler{service: service}
}

func (h *PermissionHandler) List(c *gin.Context) {
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

	sharedResponse.Paginated(c, "permissions fetched successfully", httpResponse.FromPermissionSummaries(output.Permissions), toResponseMeta(output.Meta))
}

func (h *PermissionHandler) Get(c *gin.Context) {
	output, err := h.service.Get(c.Request.Context(), c.Param("id"))
	if err != nil {
		sharedResponse.Error(c, err)
		return
	}

	sharedResponse.Success(c, http.StatusOK, "permission fetched successfully", httpResponse.FromPermissionSummary(*output))
}

func (h *PermissionHandler) Create(c *gin.Context) {
	var req httpRequest.CreatePermissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		sharedResponse.Error(c, validationError(err))
		return
	}

	output, err := h.service.Create(c.Request.Context(), dto.CreatePermissionInput{
		Name:        req.Name,
		Description: req.Description,
	})
	if err != nil {
		sharedResponse.Error(c, err)
		return
	}

	sharedResponse.Success(c, http.StatusCreated, "permission created successfully", httpResponse.FromPermissionSummary(*output))
}

func (h *PermissionHandler) Update(c *gin.Context) {
	var req httpRequest.UpdatePermissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		sharedResponse.Error(c, validationError(err))
		return
	}

	output, err := h.service.Update(c.Request.Context(), dto.UpdatePermissionInput{
		ID:          c.Param("id"),
		Name:        req.Name,
		Description: req.Description,
	})
	if err != nil {
		sharedResponse.Error(c, err)
		return
	}

	sharedResponse.Success(c, http.StatusOK, "permission updated successfully", httpResponse.FromPermissionSummary(*output))
}

func (h *PermissionHandler) Delete(c *gin.Context) {
	if err := h.service.Delete(c.Request.Context(), c.Param("id")); err != nil {
		sharedResponse.Error(c, err)
		return
	}

	sharedResponse.Success(c, http.StatusOK, "permission deleted successfully", gin.H{"deleted": true})
}
