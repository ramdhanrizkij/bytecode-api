package domain

// CreateRoleRequest is the payload for creating a new role.
// Pure Go struct — no Fiber or HTTP framework import.
type CreateRoleRequest struct {
	Name        string `json:"name"        validate:"required,min=2,max=50" example:"manager"`
	Description string `json:"description" validate:"max=255" example:"Can manage operational resources"`
	GuardName   string `json:"guard_name"  validate:"max=50" example:"api"`
}

// UpdateRoleRequest is the payload for updating an existing role.
type UpdateRoleRequest struct {
	Name        string `json:"name"        validate:"required,min=2,max=50" example:"manager"`
	Description string `json:"description" validate:"max=255" example:"Can manage operational resources"`
	GuardName   string `json:"guard_name"  validate:"max=50" example:"api"`
}

// RoleResponse is the serialised form of a role returned to clients.
type RoleResponse struct {
	ID          string               `json:"id" example:"018f7606-a3f7-7c40-8e4b-2d47c6e04c8d"`
	Name        string               `json:"name" example:"manager"`
	Description string               `json:"description" example:"Can manage operational resources"`
	GuardName   string               `json:"guard_name" example:"api"`
	Permissions []PermissionResponse `json:"permissions,omitempty"`
	CreatedAt   string               `json:"created_at" example:"2026-05-13T10:00:00+07:00"`
	UpdatedAt   string               `json:"updated_at" example:"2026-05-13T10:00:00+07:00"`
}

// PermissionResponse is the minimal permission projection embedded in RoleResponse.
type PermissionResponse struct {
	ID   string `json:"id" example:"018f7606-a3f7-7c40-8e4b-2d47c6e04c8d"`
	Name string `json:"name" example:"users.view"`
}

// AssignPermissionsRequest carries the list of permission UUIDs to assign or remove.
type AssignPermissionsRequest struct {
	PermissionIDs []string `json:"permission_ids" validate:"required,min=1,dive,uuid" example:"018f7606-a3f7-7c40-8e4b-2d47c6e04c8d"`
}
