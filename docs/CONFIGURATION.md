# Configuration

## Application Configuration

Configuration lives in `internal/core/config/config.go`. `LoadConfig` loads `.env` with `godotenv.Load()` when present, then parses environment variables with `caarlos0/env`.

## Runtime Configuration

| Group | Struct |
| --- | --- |
| App | `AppConfig` |
| Swagger | `SwaggerConfig` |
| Database | `DBConfig` |
| JWT | `JWTConfig` |
| Redis | `RedisConfig` |
| Storage | `StorageConfig` |
| Logging | `LogConfig` |

## Secrets

Secrets are read from environment variables:

- `DB_PASSWORD`
- `JWT_SECRET`
- `REDIS_PASSWORD`
- `MINIO_ACCESS_KEY`
- `MINIO_SECRET_KEY`
- `SWAGGER_PASSWORD`

## Environment-Specific Behavior

- `APP_ENV=production` makes GORM logging silent.
- `APP_ENV=production` rejects the default `JWT_SECRET`.
- Swagger in production requires `SWAGGER_ENABLED=true` and both basic auth credentials.
- Non-production Swagger is exposed without basic auth when enabled.

## Feature Flags

| Variable | Effect |
| --- | --- |
| `SWAGGER_ENABLED` | Registers or disables Swagger UI |
| `REDIS_ENABLED` | Enables Redis-backed cache or no-op cache |
| `STORAGE_PROVIDER` | Selects `local` or `minio` |
| `COMPOSE_PROFILES` | Controls optional Docker Compose Redis/MinIO services |
