-- UP: Seed default permissions and assign them to superadmin role
INSERT INTO permissions (name, description) VALUES
    ('roles.view', 'View roles list and details'),
    ('roles.create', 'Create new roles'),
    ('roles.edit', 'Update existing roles'),
    ('roles.delete', 'Delete roles'),
    ('roles.assign-permission', 'Assign permissions to roles'),
    ('roles.remove-permission', 'Remove permissions from roles'),
    ('permissions.view', 'View permissions list and details'),
    ('permissions.create', 'Create new permissions'),
    ('permissions.edit', 'Update existing permissions'),
    ('permissions.delete', 'Delete permissions'),
    ('users.view', 'View users list and details'),
    ('users.create', 'Create new users'),
    ('users.edit', 'Update existing users'),
    ('users.delete', 'Delete users')
ON CONFLICT (name) DO NOTHING;

-- Assign all permissions to superadmin role
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r, permissions p
WHERE r.name = 'superadmin'
ON CONFLICT DO NOTHING;
