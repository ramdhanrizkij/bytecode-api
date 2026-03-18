package gorm

import (
	"context"
	"errors"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/ramdhanrizki/bytecode-api/internal/identity/domain/entity"
	identityRepo "github.com/ramdhanrizki/bytecode-api/internal/identity/domain/repository"
)

type Repositories struct {
	db *gorm.DB
}

func NewRepositories(db *gorm.DB) *Repositories {
	return &Repositories{db: db}
}

func (r *Repositories) Users() identityRepo.UserRepository {
	return &UserRepository{db: r.db}
}

func (r *Repositories) Roles() identityRepo.RoleRepository {
	return &RoleRepository{db: r.db}
}

func (r *Repositories) Permissions() identityRepo.PermissionRepository {
	return &PermissionRepository{db: r.db}
}

func (r *Repositories) UserRoles() identityRepo.UserRoleRepository {
	return &UserRoleRepository{db: r.db}
}

func (r *Repositories) RolePermissions() identityRepo.RolePermissionRepository {
	return &RolePermissionRepository{db: r.db}
}

func (r *Repositories) RefreshTokens() identityRepo.RefreshTokenRepository {
	return &RefreshTokenRepository{db: r.db}
}

func (r *Repositories) EmailVerificationTokens() identityRepo.EmailVerificationTokenRepository {
	return &EmailVerificationTokenRepository{db: r.db}
}

type UnitOfWork struct {
	db *gorm.DB
}

func NewUnitOfWork(db *gorm.DB) *UnitOfWork {
	return &UnitOfWork{db: db}
}

func (u *UnitOfWork) Do(ctx context.Context, fn func(repos identityRepo.TxRepositories) error) error {
	return u.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return fn(NewRepositories(tx))
	})
}

type UserRepository struct {
	db *gorm.DB
}

func (r *UserRepository) Create(ctx context.Context, user *entity.User) error {
	return r.db.WithContext(ctx).Create(toUserModel(*user)).Error
}

func (r *UserRepository) Update(ctx context.Context, user *entity.User) error {
	return r.db.WithContext(ctx).Model(&UserModel{}).
		Where("id = ?", user.ID.String()).
		Updates(map[string]any{
			"full_name":         user.FullName,
			"email":             user.Email,
			"password_hash":     user.PasswordHash,
			"is_email_verified": user.IsEmailVerified,
			"is_active":         user.IsActive,
			"updated_at":        user.UpdatedAt,
		}).Error
}

func (r *UserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("user_id = ?", id.String()).Delete(&UserRoleModel{}).Error; err != nil {
			return err
		}
		if err := tx.Where("user_id = ?", id.String()).Delete(&RefreshTokenModel{}).Error; err != nil {
			return err
		}
		if err := tx.Where("user_id = ?", id.String()).Delete(&EmailVerificationTokenModel{}).Error; err != nil {
			return err
		}
		return tx.Where("id = ?", id.String()).Delete(&UserModel{}).Error
	})
}

func (r *UserRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	var model UserModel
	err := r.db.WithContext(ctx).Where("id = ?", id.String()).First(&model).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	user, err := toUserEntity(model)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) FindByIDWithRoles(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	var model UserModel
	err := r.db.WithContext(ctx).Preload("Roles").Where("id = ?", id.String()).First(&model).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	user, err := toUserEntity(model)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) FindByIDWithRolesAndPermissions(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	var model UserModel
	err := r.db.WithContext(ctx).Preload("Roles").Preload("Roles.Permissions").Where("id = ?", id.String()).First(&model).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	user, err := toUserEntity(model)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	var model UserModel
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&model).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	user, err := toUserEntity(model)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) FindByEmailWithRoles(ctx context.Context, email string) (*entity.User, error) {
	var model UserModel
	err := r.db.WithContext(ctx).Preload("Roles").Where("email = ?", email).First(&model).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	user, err := toUserEntity(model)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) List(ctx context.Context, options identityRepo.ListOptions) ([]entity.User, int, error) {
	query := r.db.WithContext(ctx).Model(&UserModel{})
	if search := strings.TrimSpace(options.Search); search != "" {
		like := "%" + search + "%"
		query = query.Where("full_name ILIKE ? OR email ILIKE ?", like, like)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var models []UserModel
	err := applyListOptions(query.Preload("Roles"), options, map[string]string{
		"full_name":  "full_name",
		"email":      "email",
		"created_at": "created_at",
		"updated_at": "updated_at",
	}, "created_at").Find(&models).Error
	if err != nil {
		return nil, 0, err
	}

	users := make([]entity.User, 0, len(models))
	for _, model := range models {
		user, err := toUserEntity(model)
		if err != nil {
			return nil, 0, err
		}
		users = append(users, *user)
	}

	return users, int(total), nil
}

type RoleRepository struct {
	db *gorm.DB
}

func (r *RoleRepository) Create(ctx context.Context, role *entity.Role) error {
	return r.db.WithContext(ctx).Create(&RoleModel{
		ID:          role.ID.String(),
		Name:        role.Name,
		Description: role.Description,
		CreatedAt:   role.CreatedAt,
		UpdatedAt:   role.UpdatedAt,
	}).Error
}

func (r *RoleRepository) Update(ctx context.Context, role *entity.Role) error {
	return r.db.WithContext(ctx).Model(&RoleModel{}).
		Where("id = ?", role.ID.String()).
		Updates(map[string]any{
			"name":        role.Name,
			"description": role.Description,
			"updated_at":  role.UpdatedAt,
		}).Error
}

func (r *RoleRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("role_id = ?", id.String()).Delete(&RolePermissionModel{}).Error; err != nil {
			return err
		}
		if err := tx.Where("role_id = ?", id.String()).Delete(&UserRoleModel{}).Error; err != nil {
			return err
		}
		return tx.Where("id = ?", id.String()).Delete(&RoleModel{}).Error
	})
}

func (r *RoleRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.Role, error) {
	var model RoleModel
	err := r.db.WithContext(ctx).Where("id = ?", id.String()).First(&model).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return toRoleEntity(model)
}

func (r *RoleRepository) FindByIDWithPermissions(ctx context.Context, id uuid.UUID) (*entity.Role, error) {
	var model RoleModel
	err := r.db.WithContext(ctx).Preload("Permissions").Where("id = ?", id.String()).First(&model).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return toRoleEntity(model)
}

func (r *RoleRepository) FindByIDs(ctx context.Context, ids []uuid.UUID) ([]entity.Role, error) {
	if len(ids) == 0 {
		return []entity.Role{}, nil
	}

	values := make([]string, 0, len(ids))
	for _, id := range ids {
		values = append(values, id.String())
	}

	var models []RoleModel
	err := r.db.WithContext(ctx).Where("id IN ?", values).Find(&models).Error
	if err != nil {
		return nil, err
	}

	roles := make([]entity.Role, 0, len(models))
	for _, model := range models {
		role, err := toRoleEntity(model)
		if err != nil {
			return nil, err
		}
		roles = append(roles, *role)
	}

	return roles, nil
}

func (r *RoleRepository) FindByName(ctx context.Context, name string) (*entity.Role, error) {
	var model RoleModel
	err := r.db.WithContext(ctx).Where("name = ?", name).First(&model).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	role, err := toRoleEntity(model)
	if err != nil {
		return nil, err
	}
	return role, nil
}

func (r *RoleRepository) List(ctx context.Context, options identityRepo.ListOptions) ([]entity.Role, int, error) {
	query := r.db.WithContext(ctx).Model(&RoleModel{})
	if search := strings.TrimSpace(options.Search); search != "" {
		like := "%" + search + "%"
		query = query.Where("name ILIKE ? OR description ILIKE ?", like, like)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var models []RoleModel
	err := applyListOptions(query.Preload("Permissions"), options, map[string]string{
		"name":       "name",
		"created_at": "created_at",
		"updated_at": "updated_at",
	}, "name").Find(&models).Error
	if err != nil {
		return nil, 0, err
	}

	roles := make([]entity.Role, 0, len(models))
	for _, model := range models {
		role, err := toRoleEntity(model)
		if err != nil {
			return nil, 0, err
		}
		roles = append(roles, *role)
	}

	return roles, int(total), nil
}

type PermissionRepository struct {
	db *gorm.DB
}

func (r *PermissionRepository) Create(ctx context.Context, permission *entity.Permission) error {
	return r.db.WithContext(ctx).Create(toPermissionModel(*permission)).Error
}

func (r *PermissionRepository) Update(ctx context.Context, permission *entity.Permission) error {
	return r.db.WithContext(ctx).Model(&PermissionModel{}).
		Where("id = ?", permission.ID.String()).
		Updates(map[string]any{
			"name":        permission.Name,
			"description": permission.Description,
			"updated_at":  permission.UpdatedAt,
		}).Error
}

func (r *PermissionRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("permission_id = ?", id.String()).Delete(&RolePermissionModel{}).Error; err != nil {
			return err
		}
		return tx.Where("id = ?", id.String()).Delete(&PermissionModel{}).Error
	})
}

func (r *PermissionRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.Permission, error) {
	var model PermissionModel
	err := r.db.WithContext(ctx).Where("id = ?", id.String()).First(&model).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return toPermissionEntity(model)
}

func (r *PermissionRepository) FindByIDs(ctx context.Context, ids []uuid.UUID) ([]entity.Permission, error) {
	if len(ids) == 0 {
		return []entity.Permission{}, nil
	}

	values := make([]string, 0, len(ids))
	for _, id := range ids {
		values = append(values, id.String())
	}

	var models []PermissionModel
	err := r.db.WithContext(ctx).Where("id IN ?", values).Find(&models).Error
	if err != nil {
		return nil, err
	}

	permissions := make([]entity.Permission, 0, len(models))
	for _, model := range models {
		permission, err := toPermissionEntity(model)
		if err != nil {
			return nil, err
		}
		permissions = append(permissions, *permission)
	}

	return permissions, nil
}

func (r *PermissionRepository) FindByName(ctx context.Context, name string) (*entity.Permission, error) {
	var model PermissionModel
	err := r.db.WithContext(ctx).Where("name = ?", name).First(&model).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return toPermissionEntity(model)
}

func (r *PermissionRepository) List(ctx context.Context, options identityRepo.ListOptions) ([]entity.Permission, int, error) {
	query := r.db.WithContext(ctx).Model(&PermissionModel{})
	if search := strings.TrimSpace(options.Search); search != "" {
		like := "%" + search + "%"
		query = query.Where("name ILIKE ? OR description ILIKE ?", like, like)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var models []PermissionModel
	err := applyListOptions(query, options, map[string]string{
		"name":       "name",
		"created_at": "created_at",
		"updated_at": "updated_at",
	}, "name").Find(&models).Error
	if err != nil {
		return nil, 0, err
	}

	permissions := make([]entity.Permission, 0, len(models))
	for _, model := range models {
		permission, err := toPermissionEntity(model)
		if err != nil {
			return nil, 0, err
		}
		permissions = append(permissions, *permission)
	}

	return permissions, int(total), nil
}

type UserRoleRepository struct {
	db *gorm.DB
}

func (r *UserRoleRepository) AssignRole(ctx context.Context, userID, roleID uuid.UUID) error {
	model := UserRoleModel{
		UserID: userID.String(),
		RoleID: roleID.String(),
	}
	return r.db.WithContext(ctx).Create(&model).Error
}

func (r *UserRoleRepository) ReplaceRoles(ctx context.Context, userID uuid.UUID, roleIDs []uuid.UUID) error {
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID.String()).Delete(&UserRoleModel{}).Error; err != nil {
		return err
	}
	if len(roleIDs) == 0 {
		return nil
	}

	models := make([]UserRoleModel, 0, len(roleIDs))
	for _, roleID := range roleIDs {
		models = append(models, UserRoleModel{UserID: userID.String(), RoleID: roleID.String()})
	}

	return r.db.WithContext(ctx).Clauses(clause.OnConflict{DoNothing: true}).Create(&models).Error
}

type RolePermissionRepository struct {
	db *gorm.DB
}

func (r *RolePermissionRepository) ReplacePermissions(ctx context.Context, roleID uuid.UUID, permissionIDs []uuid.UUID) error {
	if err := r.db.WithContext(ctx).Where("role_id = ?", roleID.String()).Delete(&RolePermissionModel{}).Error; err != nil {
		return err
	}
	if len(permissionIDs) == 0 {
		return nil
	}

	models := make([]RolePermissionModel, 0, len(permissionIDs))
	for _, permissionID := range permissionIDs {
		models = append(models, RolePermissionModel{RoleID: roleID.String(), PermissionID: permissionID.String()})
	}

	return r.db.WithContext(ctx).Clauses(clause.OnConflict{DoNothing: true}).Create(&models).Error
}

type RefreshTokenRepository struct {
	db *gorm.DB
}

func (r *RefreshTokenRepository) Create(ctx context.Context, token *entity.RefreshToken) error {
	return r.db.WithContext(ctx).Create(toRefreshTokenModel(*token)).Error
}

func (r *RefreshTokenRepository) FindByToken(ctx context.Context, token string) (*entity.RefreshToken, error) {
	var model RefreshTokenModel
	err := r.db.WithContext(ctx).Where("token = ?", token).First(&model).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return toRefreshTokenEntity(model)
}

func (r *RefreshTokenRepository) Update(ctx context.Context, token *entity.RefreshToken) error {
	return r.db.WithContext(ctx).Model(&RefreshTokenModel{}).
		Where("id = ?", token.ID.String()).
		Updates(map[string]any{
			"user_id":    token.UserID.String(),
			"token":      token.Token,
			"expires_at": token.ExpiresAt,
			"revoked_at": token.RevokedAt,
		}).Error
}

type EmailVerificationTokenRepository struct {
	db *gorm.DB
}

func (r *EmailVerificationTokenRepository) Create(ctx context.Context, token *entity.EmailVerificationToken) error {
	return r.db.WithContext(ctx).Create(toEmailVerificationTokenModel(*token)).Error
}

func (r *EmailVerificationTokenRepository) FindByToken(ctx context.Context, token string) (*entity.EmailVerificationToken, error) {
	var model EmailVerificationTokenModel
	err := r.db.WithContext(ctx).Where("token = ?", token).First(&model).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return toEmailVerificationTokenEntity(model)
}

func (r *EmailVerificationTokenRepository) Update(ctx context.Context, token *entity.EmailVerificationToken) error {
	return r.db.WithContext(ctx).Model(&EmailVerificationTokenModel{}).
		Where("id = ?", token.ID.String()).
		Updates(map[string]any{
			"user_id":    token.UserID.String(),
			"token":      token.Token,
			"expires_at": token.ExpiresAt,
			"used_at":    token.UsedAt,
		}).Error
}

func applyListOptions(query *gorm.DB, options identityRepo.ListOptions, allowedSorts map[string]string, defaultSort string) *gorm.DB {
	sortColumn := allowedSorts[defaultSort]
	if candidate, ok := allowedSorts[strings.ToLower(strings.TrimSpace(options.Sort))]; ok {
		sortColumn = candidate
	}
	if sortColumn == "" {
		sortColumn = defaultSort
	}

	order := "asc"
	if strings.EqualFold(strings.TrimSpace(options.Order), "desc") {
		order = "desc"
	}

	page := options.Page
	if page <= 0 {
		page = 1
	}
	limit := options.Limit
	if limit <= 0 {
		limit = 10
	}

	return query.Order(sortColumn + " " + order).Offset((page - 1) * limit).Limit(limit)
}
