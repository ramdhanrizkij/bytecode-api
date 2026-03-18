package response

import "github.com/ramdhanrizki/bytecode-api/internal/identity/application/dto"

type PermissionDetailResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type RoleDetailResponse struct {
	ID          string                     `json:"id"`
	Name        string                     `json:"name"`
	Description string                     `json:"description"`
	Permissions []PermissionDetailResponse `json:"permissions"`
}

func FromPermissionSummary(permission dto.PermissionSummary) PermissionDetailResponse {
	return PermissionDetailResponse{
		ID:          permission.ID,
		Name:        permission.Name,
		Description: permission.Description,
	}
}

func FromPermissionSummaries(permissions []dto.PermissionSummary) []PermissionDetailResponse {
	items := make([]PermissionDetailResponse, 0, len(permissions))
	for _, permission := range permissions {
		items = append(items, FromPermissionSummary(permission))
	}
	return items
}

func FromRoleSummary(role dto.RoleSummary) RoleDetailResponse {
	return RoleDetailResponse{
		ID:          role.ID,
		Name:        role.Name,
		Description: role.Description,
		Permissions: FromPermissionSummaries(role.Permissions),
	}
}

func FromRoleSummaries(roles []dto.RoleSummary) []RoleDetailResponse {
	items := make([]RoleDetailResponse, 0, len(roles))
	for _, role := range roles {
		items = append(items, FromRoleSummary(role))
	}
	return items
}

func FromUserSummaries(users []dto.UserSummary) []UserResponse {
	items := make([]UserResponse, 0, len(users))
	for _, user := range users {
		items = append(items, FromUserSummary(user))
	}
	return items
}
