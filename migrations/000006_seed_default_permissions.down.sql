-- DOWN: Remove seeded permissions (CASCADE will handle role_permissions)
DELETE FROM permissions WHERE name IN (
    'roles.view', 'roles.create', 'roles.edit', 'roles.delete', 
    'roles.assign-permission', 'roles.remove-permission',
    'permissions.view', 'permissions.create', 'permissions.edit', 'permissions.delete',
    'users.view', 'users.create', 'users.edit', 'users.delete'
);
