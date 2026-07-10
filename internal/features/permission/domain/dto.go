package domain

// CreatePermissionRequest is the payload for creating a new permission.
// Pure Go struct — no Fiber or HTTP framework import.
type CreatePermissionRequest struct {
	Name        string `json:"name"        validate:"required,min=2,max=100" example:"users.view"`
	Description string `json:"description" validate:"max=255" example:"Can view users"`
	GuardName   string `json:"guard_name"  validate:"max=50" example:"api"`
}

// UpdatePermissionRequest is the payload for updating an existing permission.
type UpdatePermissionRequest struct {
	Name        string `json:"name"        validate:"required,min=2,max=100" example:"users.view"`
	Description string `json:"description" validate:"max=255" example:"Can view users"`
	GuardName   string `json:"guard_name"  validate:"max=50" example:"api"`
}

// PermissionDetailResponse is the full serialised form of a permission.
type PermissionDetailResponse struct {
	ID          string `json:"id" example:"018f7606-a3f7-7c40-8e4b-2d47c6e04c8d"`
	Name        string `json:"name" example:"users.view"`
	Description string `json:"description" example:"Can view users"`
	GuardName   string `json:"guard_name" example:"api"`
	CreatedAt   string `json:"created_at" example:"2026-05-13T10:00:00+07:00"`
	UpdatedAt   string `json:"updated_at" example:"2026-05-13T10:00:00+07:00"`
}
