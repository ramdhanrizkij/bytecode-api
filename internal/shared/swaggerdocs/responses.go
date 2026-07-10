package swaggerdocs

import (
	authDomain "github.com/ramdhanrizkij/bytecode-api/internal/features/auth/domain"
	permDomain "github.com/ramdhanrizkij/bytecode-api/internal/features/permission/domain"
	roleDomain "github.com/ramdhanrizkij/bytecode-api/internal/features/role/domain"
	userDomain "github.com/ramdhanrizkij/bytecode-api/internal/features/user/domain"
	"github.com/ramdhanrizkij/bytecode-api/internal/shared/response"
	"github.com/ramdhanrizkij/bytecode-api/internal/shared/validator"
)

// HealthData describes the service health payload.
type HealthData struct {
	Status      string `json:"status" example:"ok"`
	Service     string `json:"service" example:"bytecode-api"`
	Environment string `json:"environment" example:"development"`
	Database    string `json:"database" example:"up"`
	Cache       string `json:"cache" example:"disabled"`
	Storage     string `json:"storage" example:"local"`
}

// HealthResponse is the health endpoint response.
type HealthResponse struct {
	Meta response.Meta `json:"meta"`
	Data HealthData    `json:"data"`
}

// ErrorResponse is the standard error response.
type ErrorResponse struct {
	Meta response.Meta `json:"meta"`
}

// ValidationErrorResponse is returned when request validation fails.
type ValidationErrorResponse struct {
	Meta   response.Meta                      `json:"meta"`
	Errors []*validator.ValidationErrorDetail `json:"errors"`
}

// AuthResponse is returned by register and login.
type AuthResponse struct {
	Meta response.Meta           `json:"meta"`
	Data authDomain.AuthResponse `json:"data"`
}

// TokenResponse is returned after refreshing tokens.
type TokenResponse struct {
	Meta response.Meta            `json:"meta"`
	Data authDomain.TokenResponse `json:"data"`
}

// UserResponse wraps a single user.
type UserResponse struct {
	Meta response.Meta                 `json:"meta"`
	Data userDomain.UserDetailResponse `json:"data"`
}

// UserListResponse wraps paginated users.
type UserListResponse struct {
	Meta       response.Meta                   `json:"meta"`
	Data       []userDomain.UserDetailResponse `json:"data"`
	Pagination *response.PaginationMeta        `json:"pagination,omitempty"`
}

// UserPermissionListResponse wraps the current user's permission names.
type UserPermissionListResponse struct {
	Meta response.Meta `json:"meta"`
	Data []string      `json:"data" example:"users.view,roles.view"`
}

// RoleResponse wraps a single role.
type RoleResponse struct {
	Meta response.Meta           `json:"meta"`
	Data roleDomain.RoleResponse `json:"data"`
}

// RoleListResponse wraps paginated roles.
type RoleListResponse struct {
	Meta       response.Meta             `json:"meta"`
	Data       []roleDomain.RoleResponse `json:"data"`
	Pagination *response.PaginationMeta  `json:"pagination,omitempty"`
}

// PermissionResponse wraps a single permission.
type PermissionResponse struct {
	Meta response.Meta                       `json:"meta"`
	Data permDomain.PermissionDetailResponse `json:"data"`
}

// PermissionListResponse wraps paginated permissions.
type PermissionListResponse struct {
	Meta       response.Meta                         `json:"meta"`
	Data       []permDomain.PermissionDetailResponse `json:"data"`
	Pagination *response.PaginationMeta              `json:"pagination,omitempty"`
}
