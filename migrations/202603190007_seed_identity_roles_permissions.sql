-- +goose Up
BEGIN;

INSERT INTO roles (id, name, description, created_at, updated_at)
VALUES
    ('00000000-0000-0000-0000-000000000001', 'admin', 'Administrator role with full access', NOW(), NOW()),
    ('00000000-0000-0000-0000-000000000002', 'user', 'Default user role', NOW(), NOW())
ON CONFLICT (name) DO UPDATE
SET description = EXCLUDED.description,
    updated_at = NOW();

INSERT INTO permissions (id, name, description, created_at, updated_at)
VALUES
    ('00000000-0000-0000-0000-000000001001', 'users.read', 'Read users', NOW(), NOW()),
    ('00000000-0000-0000-0000-000000001002', 'users.create', 'Create users', NOW(), NOW()),
    ('00000000-0000-0000-0000-000000001003', 'users.update', 'Update users', NOW(), NOW()),
    ('00000000-0000-0000-0000-000000001004', 'users.delete', 'Delete users', NOW(), NOW()),
    ('00000000-0000-0000-0000-000000001005', 'roles.read', 'Read roles', NOW(), NOW()),
    ('00000000-0000-0000-0000-000000001006', 'roles.create', 'Create roles', NOW(), NOW()),
    ('00000000-0000-0000-0000-000000001007', 'roles.update', 'Update roles', NOW(), NOW()),
    ('00000000-0000-0000-0000-000000001008', 'roles.delete', 'Delete roles', NOW(), NOW()),
    ('00000000-0000-0000-0000-000000001009', 'permissions.read', 'Read permissions', NOW(), NOW()),
    ('00000000-0000-0000-0000-000000001010', 'permissions.create', 'Create permissions', NOW(), NOW()),
    ('00000000-0000-0000-0000-000000001011', 'permissions.update', 'Update permissions', NOW(), NOW()),
    ('00000000-0000-0000-0000-000000001012', 'permissions.delete', 'Delete permissions', NOW(), NOW()),
    ('00000000-0000-0000-0000-000000001013', 'profile.read', 'Read current profile', NOW(), NOW()),
    ('00000000-0000-0000-0000-000000001014', 'profile.update', 'Update current profile', NOW(), NOW())
ON CONFLICT (name) DO UPDATE
SET description = EXCLUDED.description,
    updated_at = NOW();

INSERT INTO role_permissions (role_id, permission_id)
SELECT '00000000-0000-0000-0000-000000000001', id
FROM permissions
WHERE name IN (
    'users.read', 'users.create', 'users.update', 'users.delete',
    'roles.read', 'roles.create', 'roles.update', 'roles.delete',
    'permissions.read', 'permissions.create', 'permissions.update', 'permissions.delete',
    'profile.read', 'profile.update'
)
ON CONFLICT (role_id, permission_id) DO NOTHING;

INSERT INTO role_permissions (role_id, permission_id)
SELECT '00000000-0000-0000-0000-000000000002', id
FROM permissions
WHERE name IN ('profile.read', 'profile.update')
ON CONFLICT (role_id, permission_id) DO NOTHING;

COMMIT;

-- +goose Down
BEGIN;

DELETE FROM role_permissions
WHERE role_id IN ('00000000-0000-0000-0000-000000000001', '00000000-0000-0000-0000-000000000002');

DELETE FROM permissions
WHERE id IN (
    '00000000-0000-0000-0000-000000001001', '00000000-0000-0000-0000-000000001002',
    '00000000-0000-0000-0000-000000001003', '00000000-0000-0000-0000-000000001004',
    '00000000-0000-0000-0000-000000001005', '00000000-0000-0000-0000-000000001006',
    '00000000-0000-0000-0000-000000001007', '00000000-0000-0000-0000-000000001008',
    '00000000-0000-0000-0000-000000001009', '00000000-0000-0000-0000-000000001010',
    '00000000-0000-0000-0000-000000001011', '00000000-0000-0000-0000-000000001012',
    '00000000-0000-0000-0000-000000001013', '00000000-0000-0000-0000-000000001014'
);

DELETE FROM roles
WHERE id IN ('00000000-0000-0000-0000-000000000001', '00000000-0000-0000-0000-000000000002');

COMMIT;