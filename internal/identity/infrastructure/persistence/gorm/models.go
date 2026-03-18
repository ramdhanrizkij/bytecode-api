package gorm

import "time"

type UserModel struct {
	ID              string      `gorm:"column:id;type:uuid;primaryKey"`
	FullName        string      `gorm:"column:full_name"`
	Email           string      `gorm:"column:email"`
	PasswordHash    string      `gorm:"column:password_hash"`
	IsEmailVerified bool        `gorm:"column:is_email_verified"`
	IsActive        bool        `gorm:"column:is_active"`
	Roles           []RoleModel `gorm:"many2many:user_roles;joinForeignKey:UserID;JoinReferences:RoleID"`
	CreatedAt       time.Time   `gorm:"column:created_at"`
	UpdatedAt       time.Time   `gorm:"column:updated_at"`
}

func (UserModel) TableName() string {
	return "users"
}

type RoleModel struct {
	ID          string            `gorm:"column:id;type:uuid;primaryKey"`
	Name        string            `gorm:"column:name"`
	Description string            `gorm:"column:description"`
	Permissions []PermissionModel `gorm:"many2many:role_permissions;joinForeignKey:RoleID;JoinReferences:PermissionID"`
	CreatedAt   time.Time         `gorm:"column:created_at"`
	UpdatedAt   time.Time         `gorm:"column:updated_at"`
}

func (RoleModel) TableName() string {
	return "roles"
}

type PermissionModel struct {
	ID          string    `gorm:"column:id;type:uuid;primaryKey"`
	Name        string    `gorm:"column:name"`
	Description string    `gorm:"column:description"`
	CreatedAt   time.Time `gorm:"column:created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at"`
}

func (PermissionModel) TableName() string {
	return "permissions"
}

type RolePermissionModel struct {
	RoleID       string `gorm:"column:role_id;type:uuid;primaryKey"`
	PermissionID string `gorm:"column:permission_id;type:uuid;primaryKey"`
}

func (RolePermissionModel) TableName() string {
	return "role_permissions"
}

type UserRoleModel struct {
	UserID string `gorm:"column:user_id;type:uuid;primaryKey"`
	RoleID string `gorm:"column:role_id;type:uuid;primaryKey"`
}

func (UserRoleModel) TableName() string {
	return "user_roles"
}

type RefreshTokenModel struct {
	ID        string     `gorm:"column:id;type:uuid;primaryKey"`
	UserID    string     `gorm:"column:user_id;type:uuid"`
	Token     string     `gorm:"column:token"`
	ExpiresAt time.Time  `gorm:"column:expires_at"`
	RevokedAt *time.Time `gorm:"column:revoked_at"`
	CreatedAt time.Time  `gorm:"column:created_at"`
}

func (RefreshTokenModel) TableName() string {
	return "refresh_tokens"
}

type EmailVerificationTokenModel struct {
	ID        string     `gorm:"column:id;type:uuid;primaryKey"`
	UserID    string     `gorm:"column:user_id;type:uuid"`
	Token     string     `gorm:"column:token"`
	ExpiresAt time.Time  `gorm:"column:expires_at"`
	UsedAt    *time.Time `gorm:"column:used_at"`
	CreatedAt time.Time  `gorm:"column:created_at"`
}

func (EmailVerificationTokenModel) TableName() string {
	return "email_verification_tokens"
}
