# Troubleshooting

## Database Connection Fails

### Root Cause

PostgreSQL is not running, environment variables are wrong, or `DB_SSLMODE` does not match the server.

### Solutions

```bash
docker compose up -d postgres
make migrate-up
```

Check `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`, and `DB_NAME`.

## Production Startup Fails With JWT Secret Error

### Root Cause

`APP_ENV=production` and `JWT_SECRET=your-super-secret-key`.

### Solution

Set a strong non-default `JWT_SECRET`.

## Swagger Missing In Production

### Root Cause

Production Swagger requires `SWAGGER_ENABLED=true`, `SWAGGER_USERNAME`, and `SWAGGER_PASSWORD`. If credentials are missing, the route is not registered.

### Solution

Set both credentials or disable Swagger intentionally.

## Redis Connection Fails

### Root Cause

`REDIS_ENABLED=true` but Redis is unreachable.

### Solution

```bash
COMPOSE_PROFILES=redis docker compose up -d redis
```

Or set `REDIS_ENABLED=false`.

## MinIO Startup Fails

### Root Cause

`STORAGE_PROVIDER=minio` but MinIO endpoint, credentials, or buckets are unavailable.

### Solution

```bash
COMPOSE_PROFILES=minio docker compose up -d minio minio-init
```

Verify `MINIO_ENDPOINT`, `MINIO_ACCESS_KEY`, and `MINIO_SECRET_KEY`.

## Integration Tests Fail To Start

### Root Cause

Docker is unavailable or Testcontainers cannot start PostgreSQL.

### Solution

Start Docker, then run:

```bash
make test-integration
```

## 403 On Protected Endpoint

### Root Cause

The authenticated user's role does not have the required permission, or JWT `role_name` no longer matches a database role.

### Debugging Steps

1. Decode the JWT and inspect `role_name`.
2. Check `role_permissions` for that role.
3. Use a `superadmin` user to bypass permission checks for diagnostics.

## Useful Commands

```bash
make run
make run-worker
make test
make migrate-refresh
make swagger
docker compose ps
docker compose logs postgres
```
