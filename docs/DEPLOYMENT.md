# Deployment

## Docker

The Dockerfile is a multi-stage build:

1. Build stage: `golang:1.26.1-alpine`.
2. Installs `git`, `gcc`, and `musl-dev`.
3. Downloads modules.
4. Builds `cmd/api/main.go` to `bin/api`.
5. Runtime stage: `alpine:3.19`.
6. Copies the API binary and `.env.example` as `.env`.

```bash
docker build -t bytecode-api .
```

## Docker Compose

`docker-compose.yml` defines:

- PostgreSQL `17-alpine`.
- Redis `8-alpine` behind the `redis` profile.
- MinIO behind the `minio` profile.
- MinIO bucket initialization with `minio/mc`.

```bash
docker compose up -d postgres
COMPOSE_PROFILES=redis,minio docker compose up -d
```

## CI/CD

Not present in the analyzed codebase. No workflow files were found.

## Build Pipeline

Implemented local commands:

```bash
make build
make test
make migrate-up
make swagger
```

## Release Process

Not present in the analyzed codebase.

## Rollback Strategy

Database rollback is available one migration at a time:

```bash
make migrate-down
```

Application rollback process is not present in the analyzed codebase.

## Scaling

The API process is stateless except for in-memory RBAC permission cache and in-process worker queue. Horizontal API scaling is possible at the application level, but queued welcome-email jobs are process-local and not durable.

The worker process is separate. Multiple worker processes would duplicate scheduled jobs because no distributed lock is implemented.

## Environment Promotion

Not present in the analyzed codebase.
