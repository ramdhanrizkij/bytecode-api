package request

type CreateAdminUserRequest struct {
	FullName        string   `json:"full_name" binding:"required,min=3,max=120"`
	Email           string   `json:"email" binding:"required,email,max=255"`
	Password        string   `json:"password" binding:"required,min=8,max=72"`
	IsEmailVerified *bool    `json:"is_email_verified"`
	IsActive        *bool    `json:"is_active"`
	RoleIDs         []string `json:"role_ids"`
}

type UpdateAdminUserRequest struct {
	FullName        string `json:"full_name" binding:"required,min=3,max=120"`
	Email           string `json:"email" binding:"required,email,max=255"`
	IsEmailVerified bool   `json:"is_email_verified"`
	IsActive        bool   `json:"is_active"`
}

type AssignUserRolesRequest struct {
	RoleIDs []string `json:"role_ids" binding:"required"`
}

type CreateRoleRequest struct {
	Name        string `json:"name" binding:"required,min=2,max=100"`
	Description string `json:"description" binding:"max=255"`
}

type UpdateRoleRequest struct {
	Name        string `json:"name" binding:"required,min=2,max=100"`
	Description string `json:"description" binding:"max=255"`
}

type AssignRolePermissionsRequest struct {
	PermissionIDs []string `json:"permission_ids"`
}

type CreatePermissionRequest struct {
	Name        string `json:"name" binding:"required,min=3,max=150"`
	Description string `json:"description" binding:"max=255"`
}

type UpdatePermissionRequest struct {
	Name        string `json:"name" binding:"required,min=3,max=150"`
	Description string `json:"description" binding:"max=255"`
}
