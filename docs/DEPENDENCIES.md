# Dependencies

## Runtime

| Dependency | Why It Exists |
| --- | --- |
| `github.com/gofiber/fiber/v3` | HTTP routing and middleware |
| `github.com/gofiber/contrib/v3/swaggo` | Swagger UI handler |
| `gorm.io/gorm` | ORM for PostgreSQL access |
| `gorm.io/driver/postgres` | GORM PostgreSQL driver |
| `github.com/go-gormigrate/gormigrate/v2` | Versioned Go migration execution and rollback through GORM |
| `github.com/golang-jwt/jwt/v5` | JWT generation and validation |
| `golang.org/x/crypto` | bcrypt password hashing |
| `github.com/redis/go-redis/v9` | Optional Redis cache client |
| `github.com/minio/minio-go/v7` | Optional MinIO object storage client |
| `go.uber.org/zap` | Structured logging |
| `github.com/caarlos0/env/v11` | Environment variable parsing into config structs |
| `github.com/joho/godotenv` | Optional `.env` loading |
| `github.com/google/uuid` | UUID values in models and services |
| `github.com/go-playground/validator/v10` | Request DTO validation |

## Development

| Dependency | Why It Exists |
| --- | --- |
| `github.com/swaggo/swag` | Generates Swagger docs from annotations |

## Testing

| Dependency | Why It Exists |
| --- | --- |
| `github.com/stretchr/testify` | Test assertions and mocks |
| `github.com/testcontainers/testcontainers-go` | Integration PostgreSQL containers |
| `github.com/testcontainers/testcontainers-go/modules/postgres` | PostgreSQL module for tests |

## Infrastructure

| File | Purpose |
| --- | --- |
| `Dockerfile` | Multi-stage image for API binary |
| `docker-compose.yml` | PostgreSQL, optional Redis, optional MinIO |
| `Makefile` | Build, run, test, migration, Swagger commands |

## Indirect Dependencies

Indirect dependencies are mostly transitive dependencies for Fiber, GORM, MinIO, Testcontainers, OpenAPI generation, and telemetry libraries pulled by upstream packages.
