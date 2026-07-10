package unit

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"

	"github.com/ramdhanrizkij/bytecode-api/internal/core/cache"
	"github.com/ramdhanrizkij/bytecode-api/internal/core/storage"
	permService "github.com/ramdhanrizkij/bytecode-api/internal/features/permission/service"
	roleService "github.com/ramdhanrizkij/bytecode-api/internal/features/role/service"
	userService "github.com/ramdhanrizkij/bytecode-api/internal/features/user/service"
	"github.com/ramdhanrizkij/bytecode-api/internal/model"
	apperrors "github.com/ramdhanrizkij/bytecode-api/internal/shared/errors"
	"github.com/ramdhanrizkij/bytecode-api/internal/shared/pagination"
)

func TestRoleServiceGetAllUsesCacheAfterFirstFetch(t *testing.T) {
	repo := &fakeRoleRepository{
		roles: []model.Role{
			{
				ID:          uuid.New(),
				Name:        "admin",
				Description: "Administrator",
				GuardName:   "api",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
		},
	}
	cacheClient := newFakeCache()
	svc := roleService.NewRoleService(repo, cacheClient, time.Minute, zap.NewNop())
	query := &pagination.PaginationQuery{Page: 1, PerPage: 10, Sort: "created_at", Order: "desc"}

	firstData, firstMeta, err := svc.GetAll(context.Background(), query)
	assert.NoError(t, err)
	assert.Len(t, firstData, 1)
	assert.Equal(t, 1, repo.findAllCalls)

	secondData, secondMeta, err := svc.GetAll(context.Background(), query)
	assert.NoError(t, err)
	assert.Equal(t, firstMeta, secondMeta)
	assert.Len(t, secondData, 1)
	assert.Equal(t, firstData[0].ID, secondData[0].ID)
	assert.Equal(t, firstData[0].Name, secondData[0].Name)
	assert.Equal(t, 1, repo.findAllCalls)
}

func TestPermissionServiceDeleteInvalidatesPermissionAndRoleCaches(t *testing.T) {
	permissionID := uuid.New().String()
	repo := &fakePermissionRepository{
		permission: &model.Permission{
			ID:          uuid.MustParse(permissionID),
			Name:        "users.read",
			Description: "Read users",
			GuardName:   "api",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}
	cacheClient := newFakeCache()
	svc := permService.NewPermissionService(repo, cacheClient, time.Minute, zap.NewNop())

	err := svc.Delete(context.Background(), permissionID)

	assert.NoError(t, err)
	assert.True(t, repo.deleteCalled)
	assert.Contains(t, cacheClient.deletedPrefixes, "permissions:")
	assert.Contains(t, cacheClient.deletedPrefixes, "roles:")
}

func TestUserServiceGetPermissionsUsesCacheAfterFirstFetch(t *testing.T) {
	userID := uuid.New().String()
	repo := &fakeUserRepository{
		permissions: []model.Permission{
			{Name: "users.read"},
			{Name: "users.update"},
		},
	}
	cacheClient := newFakeCache()
	svc := userService.NewUserService(repo, cacheClient, &fakeStorage{}, "uploads", time.Minute, zap.NewNop())

	firstPerms, err := svc.GetPermissions(context.Background(), userID)
	assert.NoError(t, err)
	assert.Equal(t, []string{"users.read", "users.update"}, firstPerms)
	assert.Equal(t, 1, repo.getPermissionsCalls)

	secondPerms, err := svc.GetPermissions(context.Background(), userID)
	assert.NoError(t, err)
	assert.Equal(t, firstPerms, secondPerms)
	assert.Equal(t, 1, repo.getPermissionsCalls)
}

func TestUserServiceDeleteInvalidatesUserCaches(t *testing.T) {
	userID := uuid.New().String()
	repo := &fakeUserRepository{
		user: &model.User{
			ID:        uuid.MustParse(userID),
			Name:      "Test User",
			Email:     "user@example.com",
			IsActive:  true,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}
	cacheClient := newFakeCache()
	svc := userService.NewUserService(repo, cacheClient, &fakeStorage{}, "uploads", time.Minute, zap.NewNop())

	err := svc.Delete(context.Background(), "another-user-id", userID)

	assert.NoError(t, err)
	assert.True(t, repo.deleteCalled)
	assert.Contains(t, cacheClient.deletedPrefixes, "users:")
}

type fakeCache struct {
	store           map[string][]byte
	deletedPrefixes []string
}

func newFakeCache() *fakeCache {
	return &fakeCache{store: make(map[string][]byte)}
}

func (f *fakeCache) Get(ctx context.Context, key string, dest interface{}) (bool, error) {
	value, ok := f.store[key]
	if !ok {
		return false, nil
	}

	if err := json.Unmarshal(value, dest); err != nil {
		return false, err
	}

	return true, nil
}

func (f *fakeCache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	payload, err := json.Marshal(value)
	if err != nil {
		return err
	}
	f.store[key] = payload
	return nil
}

func (f *fakeCache) Delete(ctx context.Context, keys ...string) error {
	for _, key := range keys {
		delete(f.store, key)
	}
	return nil
}

func (f *fakeCache) DeleteByPrefix(ctx context.Context, prefix string) error {
	f.deletedPrefixes = append(f.deletedPrefixes, prefix)
	for key := range f.store {
		if len(key) >= len(prefix) && key[:len(prefix)] == prefix {
			delete(f.store, key)
		}
	}
	return nil
}

func (f *fakeCache) IsEnabled() bool {
	return true
}

func (f *fakeCache) Close() error {
	return nil
}

var _ cache.Client = (*fakeCache)(nil)

type fakeStorage struct{}

func (f *fakeStorage) Put(ctx context.Context, req *storage.PutObjectRequest) (*storage.StoredObject, error) {
	return nil, nil
}

func (f *fakeStorage) Delete(ctx context.Context, bucket, key string) error {
	return nil
}

func (f *fakeStorage) URL(bucket, key string) string {
	return "https://example.test/" + bucket + "/" + key
}

func (f *fakeStorage) ProviderName() string {
	return storage.ProviderLocal
}

func (f *fakeStorage) Close() error {
	return nil
}

type fakeRoleRepository struct {
	roles        []model.Role
	findAllCalls int
}

func (f *fakeRoleRepository) FindAll(ctx context.Context, pq *pagination.PaginationQuery) ([]model.Role, int64, error) {
	f.findAllCalls++
	return f.roles, int64(len(f.roles)), nil
}

func (f *fakeRoleRepository) FindByID(ctx context.Context, id string) (*model.Role, error) {
	return nil, apperrors.ErrNotFound
}

func (f *fakeRoleRepository) FindByName(ctx context.Context, name string) (*model.Role, error) {
	return nil, apperrors.ErrNotFound
}

func (f *fakeRoleRepository) Create(ctx context.Context, role *model.Role) error {
	return nil
}

func (f *fakeRoleRepository) Update(ctx context.Context, role *model.Role) error {
	return nil
}

func (f *fakeRoleRepository) Delete(ctx context.Context, id string) error {
	return nil
}

func (f *fakeRoleRepository) AssignPermissions(ctx context.Context, roleID string, permissionIDs []string) error {
	return nil
}

func (f *fakeRoleRepository) RemovePermissions(ctx context.Context, roleID string, permissionIDs []string) error {
	return nil
}

type fakePermissionRepository struct {
	permission   *model.Permission
	deleteCalled bool
}

func (f *fakePermissionRepository) FindAll(ctx context.Context, pq *pagination.PaginationQuery) ([]model.Permission, int64, error) {
	return nil, 0, nil
}

func (f *fakePermissionRepository) FindByID(ctx context.Context, id string) (*model.Permission, error) {
	if f.permission == nil {
		return nil, apperrors.ErrNotFound
	}
	return f.permission, nil
}

func (f *fakePermissionRepository) FindByName(ctx context.Context, name string) (*model.Permission, error) {
	return nil, apperrors.ErrNotFound
}

func (f *fakePermissionRepository) FindByIDs(ctx context.Context, ids []string) ([]model.Permission, error) {
	return nil, nil
}

func (f *fakePermissionRepository) Create(ctx context.Context, permission *model.Permission) error {
	return nil
}

func (f *fakePermissionRepository) Update(ctx context.Context, permission *model.Permission) error {
	return nil
}

func (f *fakePermissionRepository) Delete(ctx context.Context, id string) error {
	f.deleteCalled = true
	return nil
}

type fakeUserRepository struct {
	user                *model.User
	permissions         []model.Permission
	deleteCalled        bool
	getPermissionsCalls int
}

func (f *fakeUserRepository) FindAll(ctx context.Context, pq *pagination.PaginationQuery) ([]model.User, int64, error) {
	if f.user == nil {
		return nil, 0, nil
	}
	return []model.User{*f.user}, 1, nil
}

func (f *fakeUserRepository) FindByID(ctx context.Context, id string) (*model.User, error) {
	if f.user == nil {
		return nil, apperrors.ErrNotFound
	}
	return f.user, nil
}

func (f *fakeUserRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	return nil, apperrors.ErrNotFound
}

func (f *fakeUserRepository) Create(ctx context.Context, user *model.User) error {
	return nil
}

func (f *fakeUserRepository) Update(ctx context.Context, user *model.User) error {
	return nil
}

func (f *fakeUserRepository) Delete(ctx context.Context, id string) error {
	f.deleteCalled = true
	return nil
}

func (f *fakeUserRepository) GetPermissions(ctx context.Context, userID string) ([]model.Permission, error) {
	f.getPermissionsCalls++
	return f.permissions, nil
}
