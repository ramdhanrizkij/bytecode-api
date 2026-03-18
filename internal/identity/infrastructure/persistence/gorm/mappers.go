package gorm

import (
	"github.com/google/uuid"

	"github.com/ramdhanrizki/bytecode-api/internal/identity/domain/entity"
)

func toUserModel(user entity.User) *UserModel {
	return &UserModel{
		ID:              user.ID.String(),
		FullName:        user.FullName,
		Email:           user.Email,
		PasswordHash:    user.PasswordHash,
		IsEmailVerified: user.IsEmailVerified,
		IsActive:        user.IsActive,
		CreatedAt:       user.CreatedAt,
		UpdatedAt:       user.UpdatedAt,
	}
}

func toUserEntity(model UserModel) (*entity.User, error) {
	id, err := uuid.Parse(model.ID)
	if err != nil {
		return nil, err
	}

	roles := make([]entity.Role, 0, len(model.Roles))
	for _, roleModel := range model.Roles {
		role, err := toRoleEntity(roleModel)
		if err != nil {
			return nil, err
		}
		roles = append(roles, *role)
	}

	return &entity.User{
		ID:              id,
		FullName:        model.FullName,
		Email:           model.Email,
		PasswordHash:    model.PasswordHash,
		IsEmailVerified: model.IsEmailVerified,
		IsActive:        model.IsActive,
		Roles:           roles,
		CreatedAt:       model.CreatedAt,
		UpdatedAt:       model.UpdatedAt,
	}, nil
}

func toRoleEntity(model RoleModel) (*entity.Role, error) {
	id, err := uuid.Parse(model.ID)
	if err != nil {
		return nil, err
	}

	permissions := make([]entity.Permission, 0, len(model.Permissions))
	for _, permissionModel := range model.Permissions {
		permission, err := toPermissionEntity(permissionModel)
		if err != nil {
			return nil, err
		}
		permissions = append(permissions, *permission)
	}

	return &entity.Role{
		ID:          id,
		Name:        model.Name,
		Description: model.Description,
		Permissions: permissions,
		CreatedAt:   model.CreatedAt,
		UpdatedAt:   model.UpdatedAt,
	}, nil
}

func toPermissionModel(permission entity.Permission) *PermissionModel {
	return &PermissionModel{
		ID:          permission.ID.String(),
		Name:        permission.Name,
		Description: permission.Description,
		CreatedAt:   permission.CreatedAt,
		UpdatedAt:   permission.UpdatedAt,
	}
}

func toPermissionEntity(model PermissionModel) (*entity.Permission, error) {
	id, err := uuid.Parse(model.ID)
	if err != nil {
		return nil, err
	}

	return &entity.Permission{
		ID:          id,
		Name:        model.Name,
		Description: model.Description,
		CreatedAt:   model.CreatedAt,
		UpdatedAt:   model.UpdatedAt,
	}, nil
}

func toRefreshTokenModel(token entity.RefreshToken) *RefreshTokenModel {
	return &RefreshTokenModel{
		ID:        token.ID.String(),
		UserID:    token.UserID.String(),
		Token:     token.Token,
		ExpiresAt: token.ExpiresAt,
		RevokedAt: token.RevokedAt,
		CreatedAt: token.CreatedAt,
	}
}

func toRefreshTokenEntity(model RefreshTokenModel) (*entity.RefreshToken, error) {
	id, err := uuid.Parse(model.ID)
	if err != nil {
		return nil, err
	}
	userID, err := uuid.Parse(model.UserID)
	if err != nil {
		return nil, err
	}

	return &entity.RefreshToken{
		ID:        id,
		UserID:    userID,
		Token:     model.Token,
		ExpiresAt: model.ExpiresAt,
		RevokedAt: model.RevokedAt,
		CreatedAt: model.CreatedAt,
	}, nil
}

func toEmailVerificationTokenModel(token entity.EmailVerificationToken) *EmailVerificationTokenModel {
	return &EmailVerificationTokenModel{
		ID:        token.ID.String(),
		UserID:    token.UserID.String(),
		Token:     token.Token,
		ExpiresAt: token.ExpiresAt,
		UsedAt:    token.UsedAt,
		CreatedAt: token.CreatedAt,
	}
}

func toEmailVerificationTokenEntity(model EmailVerificationTokenModel) (*entity.EmailVerificationToken, error) {
	id, err := uuid.Parse(model.ID)
	if err != nil {
		return nil, err
	}
	userID, err := uuid.Parse(model.UserID)
	if err != nil {
		return nil, err
	}

	return &entity.EmailVerificationToken{
		ID:        id,
		UserID:    userID,
		Token:     model.Token,
		ExpiresAt: model.ExpiresAt,
		UsedAt:    model.UsedAt,
		CreatedAt: model.CreatedAt,
	}, nil
}
