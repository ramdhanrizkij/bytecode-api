package repository

import "context"

type TxRepositories interface {
	Users() UserRepository
	Roles() RoleRepository
	Permissions() PermissionRepository
	UserRoles() UserRoleRepository
	RolePermissions() RolePermissionRepository
	RefreshTokens() RefreshTokenRepository
	EmailVerificationTokens() EmailVerificationTokenRepository
}

type UnitOfWork interface {
	Do(ctx context.Context, fn func(repos TxRepositories) error) error
}
