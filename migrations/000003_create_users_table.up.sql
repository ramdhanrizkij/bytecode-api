-- Create users table
-- role_id references roles; SET NULL if role is deleted.
CREATE TABLE users (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name       VARCHAR(100) NOT NULL,
    email      VARCHAR(100) UNIQUE NOT NULL,
    password   VARCHAR(255) NOT NULL,
    role_id    UUID         REFERENCES roles(id) ON DELETE SET NULL,
    is_active  BOOLEAN      DEFAULT true,
    created_at TIMESTAMPTZ  DEFAULT NOW(),
    updated_at TIMESTAMPTZ  DEFAULT NOW()
);
