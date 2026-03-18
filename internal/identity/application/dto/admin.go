package dto

type ListInput struct {
	Page   int
	Limit  int
	Search string
	Sort   string
	Order  string
}

type PaginationMeta struct {
	Page       int
	Limit      int
	Total      int
	TotalPages int
}

type PermissionSummary struct {
	ID          string
	Name        string
	Description string
}

type RoleSummary struct {
	ID          string
	Name        string
	Description string
	Permissions []PermissionSummary
}

type UserListOutput struct {
	Users []UserSummary
	Meta  PaginationMeta
}

type RoleListOutput struct {
	Roles []RoleSummary
	Meta  PaginationMeta
}

type PermissionListOutput struct {
	Permissions []PermissionSummary
	Meta        PaginationMeta
}

type CreateAdminUserInput struct {
	FullName        string
	Email           string
	Password        string
	IsEmailVerified *bool
	IsActive        *bool
	RoleIDs         []string
}

type UpdateAdminUserInput struct {
	ID              string
	FullName        string
	Email           string
	IsEmailVerified bool
	IsActive        bool
}

type AssignUserRolesInput struct {
	UserID  string
	RoleIDs []string
}

type CreateRoleInput struct {
	Name        string
	Description string
}

type UpdateRoleInput struct {
	ID          string
	Name        string
	Description string
}

type AssignRolePermissionsInput struct {
	RoleID        string
	PermissionIDs []string
}

type CreatePermissionInput struct {
	Name        string
	Description string
}

type UpdatePermissionInput struct {
	ID          string
	Name        string
	Description string
}
