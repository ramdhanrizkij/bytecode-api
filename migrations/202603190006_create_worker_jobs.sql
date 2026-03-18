-- +goose Up
BEGIN;

CREATE TABLE IF NOT EXISTS worker_jobs (
    id UUID PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    payload BYTEA NOT NULL,
    attempts INTEGER NOT NULL DEFAULT 0,
    max_attempts INTEGER NOT NULL DEFAULT 5,
    run_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    reserved_at TIMESTAMPTZ NULL,
    completed_at TIMESTAMPTZ NULL,
    failed_at TIMESTAMPTZ NULL,
    last_error TEXT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_worker_jobs_fetch ON worker_jobs (run_at, created_at)
WHERE completed_at IS NULL AND failed_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_worker_jobs_name ON worker_jobs (name);
CREATE INDEX IF NOT EXISTS idx_worker_jobs_reserved_at ON worker_jobs (reserved_at);

COMMIT;

-- +goose Down
BEGIN;

DROP TABLE IF EXISTS worker_jobs;

COMMIT;