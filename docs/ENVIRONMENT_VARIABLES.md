# Environment Variables

## Variables

| Variable | Required | Default | Description | Security Notes |
| --- | --- | --- | --- | --- |
| `APP_NAME` | No | `bytecode-api` | Fiber application name and health response service name | Not secret |
| `APP_PORT` | No | `8080` | HTTP listen port | Not secret |
| `APP_ENV` | No | `development` | Runtime environment | Set `production` in production |
| `SWAGGER_ENABLED` | No | `true` | Enables Swagger route | Disable or protect in production |
| `SWAGGER_USERNAME` | Production Swagger only | empty | Swagger basic auth username | Treat as sensitive |
| `SWAGGER_PASSWORD` | Production Swagger only | empty | Swagger basic auth password | Secret |
| `DB_HOST` | No | `localhost` | PostgreSQL host | Not secret |
| `DB_PORT` | No | `5432` | PostgreSQL port | Not secret |
| `DB_USER` | No | `postgres` | PostgreSQL username | Sensitive in shared logs |
| `DB_PASSWORD` | No | `secret` | PostgreSQL password | Secret |
| `DB_NAME` | No | `bytecode_api` | PostgreSQL database name | Not secret |
| `DB_SSLMODE` | No | `disable` | PostgreSQL SSL mode | Use secure mode when required |
| `JWT_SECRET` | No | `your-super-secret-key` | HMAC signing key | Must be changed in production |
| `JWT_EXPIRY_HOURS` | No | `24` | Access token lifetime | Shorter values reduce token exposure |
| `JWT_REFRESH_EXPIRY_HOURS` | No | `168` | Refresh token lifetime | Longer values increase session risk |
| `REDIS_ENABLED` | No | `false` | Enables Redis cache client | Not secret |
| `REDIS_HOST` | No | `localhost` | Redis host | Not secret |
| `REDIS_PORT` | No | `6379` | Redis port | Not secret |
| `REDIS_PASSWORD` | No | empty | Redis password | Secret when used |
| `REDIS_DB` | No | `0` | Redis logical database | Not secret |
| `REDIS_CACHE_TTL_MINUTES` | No | `5` | Cache TTL for service caches | Not secret |
| `STORAGE_PROVIDER` | No | `local` | `local` or `minio` | Not secret |
| `STORAGE_DEFAULT_BUCKET` | No | `uploads` | Default storage bucket | Not secret |
| `STORAGE_BUCKETS` | No | `uploads` | Comma-separated buckets | Not secret |
| `STORAGE_LOCAL_PATH` | No | `storage` | Local storage root | Must be writable by process |
| `STORAGE_BASE_URL` | No | `/storage` | Public local storage route prefix | Avoid exposing unintended paths |
| `MINIO_ENDPOINT` | No | `localhost:9000` | MinIO endpoint | Not secret |
| `MINIO_PUBLIC_URL` | No | `http://localhost:9000` | Public URL used in object URLs | Avoid internal-only URLs in public APIs |
| `MINIO_ACCESS_KEY` | No | `minioadmin` | MinIO access key | Secret |
| `MINIO_SECRET_KEY` | No | `minioadmin` | MinIO secret key | Secret |
| `MINIO_USE_SSL` | No | `false` | Enables HTTPS for MinIO client | Use true for TLS endpoints |
| `MINIO_REGION` | No | `us-east-1` | MinIO region | Not secret |
| `COMPOSE_PROFILES` | No | empty | Docker Compose optional services | Not read by app |
| `LOG_LEVEL` | No | `info` | `debug`, `info`, `warn`, or `error` | Debug logs can be verbose |
