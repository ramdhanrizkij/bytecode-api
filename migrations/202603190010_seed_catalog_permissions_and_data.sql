-- +goose Up
BEGIN;

INSERT INTO permissions (id, name, description, created_at, updated_at)
VALUES
    ('00000000-0000-0000-0000-000000001015', 'categories.read', 'Read categories', NOW(), NOW()),
    ('00000000-0000-0000-0000-000000001016', 'categories.create', 'Create categories', NOW(), NOW()),
    ('00000000-0000-0000-0000-000000001017', 'categories.update', 'Update categories', NOW(), NOW()),
    ('00000000-0000-0000-0000-000000001018', 'categories.delete', 'Delete categories', NOW(), NOW()),
    ('00000000-0000-0000-0000-000000001019', 'products.read', 'Read products', NOW(), NOW()),
    ('00000000-0000-0000-0000-000000001020', 'products.create', 'Create products', NOW(), NOW()),
    ('00000000-0000-0000-0000-000000001021', 'products.update', 'Update products', NOW(), NOW()),
    ('00000000-0000-0000-0000-000000001022', 'products.delete', 'Delete products', NOW(), NOW())
ON CONFLICT (name) DO UPDATE
SET description = EXCLUDED.description,
    updated_at = NOW();

INSERT INTO role_permissions (role_id, permission_id)
SELECT '00000000-0000-0000-0000-000000000001', id
FROM permissions
WHERE name IN (
    'categories.read', 'categories.create', 'categories.update', 'categories.delete',
    'products.read', 'products.create', 'products.update', 'products.delete'
)
ON CONFLICT (role_id, permission_id) DO NOTHING;

INSERT INTO role_permissions (role_id, permission_id)
SELECT '00000000-0000-0000-0000-000000000002', id
FROM permissions
WHERE name IN ('categories.read', 'products.read')
ON CONFLICT (role_id, permission_id) DO NOTHING;

INSERT INTO categories (id, name, slug, description, is_active, created_at, updated_at)
VALUES
    ('10000000-0000-0000-0000-000000000001', 'Electronics', 'electronics', 'Electronic devices and accessories', TRUE, NOW(), NOW()),
    ('10000000-0000-0000-0000-000000000002', 'Books', 'books', 'Books and reading materials', TRUE, NOW(), NOW()),
    ('10000000-0000-0000-0000-000000000003', 'Fashion', 'fashion', 'Clothing and apparel', TRUE, NOW(), NOW())
ON CONFLICT (slug) DO UPDATE
SET name = EXCLUDED.name,
    description = EXCLUDED.description,
    is_active = EXCLUDED.is_active,
    updated_at = NOW();

INSERT INTO products (id, category_id, name, slug, description, sku, price, stock, is_active, created_at, updated_at)
VALUES
    ('20000000-0000-0000-0000-000000000001', '10000000-0000-0000-0000-000000000001', 'Mechanical Keyboard', 'mechanical-keyboard', 'RGB mechanical keyboard', 'ELEC-KB-001', 850000, 25, TRUE, NOW(), NOW()),
    ('20000000-0000-0000-0000-000000000002', '10000000-0000-0000-0000-000000000001', 'Wireless Mouse', 'wireless-mouse', 'Ergonomic wireless mouse', 'ELEC-MS-001', 350000, 40, TRUE, NOW(), NOW()),
    ('20000000-0000-0000-0000-000000000003', '10000000-0000-0000-0000-000000000002', 'Clean Architecture Book', 'clean-architecture-book', 'Book by Robert C. Martin', 'BOOK-CA-001', 450000, 15, TRUE, NOW(), NOW()),
    ('20000000-0000-0000-0000-000000000004', '10000000-0000-0000-0000-000000000002', 'Domain-Driven Design Book', 'domain-driven-design-book', 'Book by Eric Evans', 'BOOK-DDD-001', 550000, 10, TRUE, NOW(), NOW()),
    ('20000000-0000-0000-0000-000000000005', '10000000-0000-0000-0000-000000000003', 'Basic T-Shirt', 'basic-t-shirt', 'Comfortable cotton t-shirt', 'FSHN-TS-001', 120000, 60, TRUE, NOW(), NOW()),
    ('20000000-0000-0000-0000-000000000006', '10000000-0000-0000-0000-000000000003', 'Hoodie', 'hoodie', 'Casual everyday hoodie', 'FSHN-HD-001', 275000, 30, TRUE, NOW(), NOW())
ON CONFLICT (sku) DO UPDATE
SET category_id = EXCLUDED.category_id,
    name = EXCLUDED.name,
    slug = EXCLUDED.slug,
    description = EXCLUDED.description,
    price = EXCLUDED.price,
    stock = EXCLUDED.stock,
    is_active = EXCLUDED.is_active,
    updated_at = NOW();

COMMIT;

-- +goose Down
BEGIN;

DELETE FROM products
WHERE id IN (
    '20000000-0000-0000-0000-000000000001',
    '20000000-0000-0000-0000-000000000002',
    '20000000-0000-0000-0000-000000000003',
    '20000000-0000-0000-0000-000000000004',
    '20000000-0000-0000-0000-000000000005',
    '20000000-0000-0000-0000-000000000006'
);

DELETE FROM categories
WHERE id IN (
    '10000000-0000-0000-0000-000000000001',
    '10000000-0000-0000-0000-000000000002',
    '10000000-0000-0000-0000-000000000003'
);

DELETE FROM role_permissions
WHERE permission_id IN (
    '00000000-0000-0000-0000-000000001015', '00000000-0000-0000-0000-000000001016',
    '00000000-0000-0000-0000-000000001017', '00000000-0000-0000-0000-000000001018',
    '00000000-0000-0000-0000-000000001019', '00000000-0000-0000-0000-000000001020',
    '00000000-0000-0000-0000-000000001021', '00000000-0000-0000-0000-000000001022'
);

DELETE FROM permissions
WHERE id IN (
    '00000000-0000-0000-0000-000000001015', '00000000-0000-0000-0000-000000001016',
    '00000000-0000-0000-0000-000000001017', '00000000-0000-0000-0000-000000001018',
    '00000000-0000-0000-0000-000000001019', '00000000-0000-0000-0000-000000001020',
    '00000000-0000-0000-0000-000000001021', '00000000-0000-0000-0000-000000001022'
);

COMMIT;