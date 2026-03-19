-- +goose Up
BEGIN;

CREATE TABLE IF NOT EXISTS categories (
    id UUID PRIMARY KEY,
    name VARCHAR(120) NOT NULL UNIQUE,
    slug VARCHAR(150) NOT NULL UNIQUE,
    description TEXT NULL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_categories_created_at ON categories (created_at DESC);
CREATE INDEX IF NOT EXISTS idx_categories_slug ON categories (slug);

COMMIT;

-- +goose Down
BEGIN;

DROP TABLE IF EXISTS categories;

COMMIT;