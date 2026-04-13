package domain

// CreatePermissionRequest is the payload for creating a new permission.
// Pure Go struct — no Fiber or HTTP framework import.
type CreatePermissionRequest struct {
	Name        string `json:"name"        validate:"required,min=2,max=100"`
	Description string `json:"description" validate:"max=255"`
	GuardName   string `json:"guard_name"  validate:"max=50"`
}

// UpdatePermissionRequest is the payload for updating an existing permission.
type UpdatePermissionRequest struct {
	Name        string `json:"name"        validate:"required,min=2,max=100"`
	Description string `json:"description" validate:"max=255"`
	GuardName   string `json:"guard_name"  validate:"max=50"`
}

// PermissionDetailResponse is the full serialised form of a permission.
type PermissionDetailResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	GuardName   string `json:"guard_name"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}
