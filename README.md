# Bytecode API

A modular Go backend boilerplate for authentication, authorization, catalog management, profile management, and asynchronous email verification.

This project uses:
- Go
- Gin
- GORM + PostgreSQL
- JWT access tokens + database-backed refresh tokens
- Zap logging
- Goose SQL migrations
- Separate API and Worker processes
- DB-backed async job queue for email verification

## Current Scope

The current implementation includes:
- category module with authenticated read and admin CRUD
- user registration
- email verification
- login
- refresh token flow
- current user profile read/update
- admin user CRUD
- product module with authenticated read and admin CRUD
- role CRUD
- permission CRUD
- assign roles to users
- assign permissions to roles
- permission-based authorization middleware
- async email verification delivery through a worker process
- health endpoints

## Architecture

The codebase follows a modular monolith structure with clear layer boundaries:
- `domain` → entities, business rules, repository contracts
- `application` → use cases / services / DTOs
- `infrastructure` → GORM repositories, queue publisher, persistence adapters
- `delivery` → HTTP handlers, middleware, request/response DTOs, worker handlers

Important top-level directories:

```text
cmd/
  api/
  worker/
configs/
internal/
  bootstrap/
  category/
  identity/
    application/
    domain/
    infrastructure/
    delivery/
      http/
      worker/
  platform/
  product/
  shared/
  worker/
migrations/
```

## Processes

### API process
Handles HTTP traffic and publishes background jobs.

Entry point:
- `cmd/api/main.go`

### Worker process
Consumes queued jobs from the `worker_jobs` table and sends verification emails asynchronously.

Entry point:
- `cmd/worker/main.go`

## Async Email Flow

Registration is non-blocking for email sending:
1. user registers via API
2. user + verification token are persisted
3. an `identity.email_verification.send` job is inserted into `worker_jobs`
4. worker polls the queue table
5. worker sends email via SMTP
6. failed jobs are retried with backoff up to the configured max attempts

## Requirements

- Go `1.24+`
- PostgreSQL
- `make` (optional but recommended)
- SMTP server for email delivery

For local email testing, MailHog/Mailpit works well.

## Environment Configuration

Copy `.env.example` to `.env` and adjust values.

Required variables:

```dotenv
APP_NAME=bytecode-api
APP_ENV=development
APP_PORT=8080
APP_BASE_URL=http://localhost:8080

DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=bytecode
DB_SSLMODE=disable

JWT_SECRET=change-me
JWT_ISSUER=bytecode-api
JWT_ACCESS_TTL_MINUTES=15
JWT_REFRESH_TTL_HOURS=168

SMTP_HOST=localhost
SMTP_PORT=1025
SMTP_USERNAME=
SMTP_PASSWORD=
SMTP_FROM=no-reply@example.com
SMTP_FROM_NAME=Bytecode API

WORKER_CONCURRENCY=5
```

Reference file:
- [.env.example](.env.example)

## Installation

Install dependencies:

```bash
go mod download
```

If you want to use Makefile-based migration commands, ensure `make` is installed on your machine.

## Database Migrations

Goose SQL migrations are stored in [migrations](migrations).

Included migrations:
- categories
- users
- roles
- permissions
- products
- user_roles
- role_permissions
- refresh_tokens
- email_verification_tokens
- worker_jobs
- default role/permission seed data
- sample category and product seed data

### Using Makefile

```bash
make install-goose
make migrate-up
make migrate-status
```

Other commands:

```bash
make migrate-down
make migrate-reset
make goose-create name=create_something
```

### Manual Goose usage

```bash
go run github.com/pressly/goose/v3/cmd/goose@v3.24.1 -dir migrations postgres "host=localhost port=5432 user=postgres password=postgres dbname=bytecode sslmode=disable TimeZone=UTC" up
```

## Run the Project

### With Makefile

Run API:

```bash
make run-api
```

Run Worker:

```bash
make run-worker
```

Build binaries:

```bash
make build
```

Or rebuild from scratch:

```bash
make build-all
```

### Without Makefile

Run API:

```bash
go run ./cmd/api
```

Run Worker:

```bash
go run ./cmd/worker
```

Build binaries:

```bash
go build -o bin/api ./cmd/api
go build -o bin/worker ./cmd/worker
```

## Health Endpoints

Public health checks:
- `GET /health`
- `GET /api/v1/health`

Example response:

```json
{
  "success": true,
  "message": "service is healthy",
  "data": {
    "status": "ok",
    "name": "bytecode-api",
    "env": "development",
    "timestamp": "2026-03-19T00:00:00Z"
  }
}
```

## Main API Endpoints

### Public
- `POST /api/v1/auth/register`
- `POST /api/v1/auth/verify-email`
- `POST /api/v1/auth/login`
- `POST /api/v1/auth/refresh`

### Authenticated user
- `GET /api/v1/profile`
- `PUT /api/v1/profile`
- `GET /api/v1/categories`
- `GET /api/v1/categories/:id`
- `GET /api/v1/products`
- `GET /api/v1/products/:id`

### Admin catalog
- `GET /api/v1/admin/categories`
- `POST /api/v1/admin/categories`
- `GET /api/v1/admin/categories/:id`
- `PUT /api/v1/admin/categories/:id`
- `DELETE /api/v1/admin/categories/:id`

- `GET /api/v1/admin/products`
- `POST /api/v1/admin/products`
- `GET /api/v1/admin/products/:id`
- `PUT /api/v1/admin/products/:id`
- `DELETE /api/v1/admin/products/:id`

### Admin
- `GET /api/v1/admin/users`
- `POST /api/v1/admin/users`
- `GET /api/v1/admin/users/:id`
- `PUT /api/v1/admin/users/:id`
- `DELETE /api/v1/admin/users/:id`
- `PUT /api/v1/admin/users/:id/roles`

- `GET /api/v1/admin/roles`
- `POST /api/v1/admin/roles`
- `GET /api/v1/admin/roles/:id`
- `PUT /api/v1/admin/roles/:id`
- `DELETE /api/v1/admin/roles/:id`
- `PUT /api/v1/admin/roles/:id/permissions`

- `GET /api/v1/admin/permissions`
- `POST /api/v1/admin/permissions`
- `GET /api/v1/admin/permissions/:id`
- `PUT /api/v1/admin/permissions/:id`
- `DELETE /api/v1/admin/permissions/:id`

## Authorization Model

Permission checks are enforced in middleware.

Seeded permissions:
- `categories.read`
- `categories.create`
- `categories.update`
- `categories.delete`
- `products.read`
- `products.create`
- `products.update`
- `products.delete`
- `users.read`
- `users.create`
- `users.update`
- `users.delete`
- `roles.read`
- `roles.create`
- `roles.update`
- `roles.delete`
- `permissions.read`
- `permissions.create`
- `permissions.update`
- `permissions.delete`
- `profile.read`
- `profile.update`

Seeded roles:
- `admin` → all seeded permissions, including catalog management
- `user` → `profile.read`, `profile.update`, `categories.read`, `products.read`

## Catalog Model

### Category
- `id`
- `name`
- `slug`
- `description`
- `is_active`
- `created_at`
- `updated_at`

### Product
- `id`
- `category_id`
- `name`
- `slug`
- `description`
- `sku`
- `price` as `BIGINT` in minor units
- `stock`
- `is_active`
- `created_at`
- `updated_at`

The product list/detail responses also include `category_name` for convenience.

## Response Format

Successful responses use a consistent envelope:

```json
{
  "success": true,
  "message": "...",
  "data": {}
}
```

Validation/error responses:

```json
{
  "success": false,
  "message": "validation failed",
  "errors": {
    "field": ["message"]
  }
}
```

Paginated responses:

```json
{
  "success": true,
  "message": "...",
  "data": [],
  "meta": {
    "page": 1,
    "limit": 10,
    "total": 100,
    "total_pages": 10
  }
}
```

## Logging

Zap is used for:
- app startup/shutdown logs
- HTTP request logging
- database initialization logs
- worker lifecycle logs
- queue processing logs
- SMTP delivery logs

## Notes

- The project currently has no automated tests.
- The `identity` module currently also contains admin-facing user/role/permission management endpoints.
- `category` and `product` are separate modules wired through the same auth and permission middleware used by `identity`.
- Schema setup is migration-driven; run migrations before starting API or worker.
- Build artifacts are written to `bin/` and ignored via `.gitignore`.
