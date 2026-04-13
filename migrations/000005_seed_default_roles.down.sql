-- Remove the default seeded roles
DELETE FROM roles WHERE name IN ('superadmin', 'admin', 'user');
