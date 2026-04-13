package domain

// CreateRoleRequest is the payload for creating a new role.
// Pure Go struct — no Fiber or HTTP framework import.
type CreateRoleRequest struct {
	Name        string `json:"name"        validate:"required,min=2,max=50"`
	Description string `json:"description" validate:"max=255"`
	GuardName   string `json:"guard_name"  validate:"max=50"`
}

// UpdateRoleRequest is the payload for updating an existing role.
type UpdateRoleRequest struct {
	Name        string `json:"name"        validate:"required,min=2,max=50"`
	Description string `json:"description" validate:"max=255"`
	GuardName   string `json:"guard_name"  validate:"max=50"`
}

// RoleResponse is the serialised form of a role returned to clients.
type RoleResponse struct {
	ID          string               `json:"id"`
	Name        string               `json:"name"`
	Description string               `json:"description"`
	GuardName   string               `json:"guard_name"`
	Permissions []PermissionResponse `json:"permissions,omitempty"`
	CreatedAt   string               `json:"created_at"`
	UpdatedAt   string               `json:"updated_at"`
}

// PermissionResponse is the minimal permission projection embedded in RoleResponse.
type PermissionResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// AssignPermissionsRequest carries the list of permission UUIDs to assign or remove.
type AssignPermissionsRequest struct {
	PermissionIDs []string `json:"permission_ids" validate:"required,min=1,dive,uuid"`
}
