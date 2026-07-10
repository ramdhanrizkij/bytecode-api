# Request Lifecycle

## HTTP Request Flow

```mermaid
sequenceDiagram
  participant Client
  participant Fiber
  participant GlobalMiddleware
  participant RouteMiddleware
  participant Handler
  participant Validator
  participant Service
  participant Repository
  participant DB
  participant Cache
  Client->>Fiber: HTTP request
  Fiber->>GlobalMiddleware: recover, CORS, request logger
  GlobalMiddleware->>RouteMiddleware: JWT and RBAC when protected
  RouteMiddleware->>Handler: matched route
  Handler->>Validator: ParseAndValidate for JSON bodies
  Handler->>Service: domain request
  Service->>Cache: read cache when implemented
  alt cache hit
    Cache-->>Service: cached response
  else cache miss
    Service->>Repository: data access
    Repository->>DB: GORM query
    DB-->>Repository: rows
    Repository-->>Service: models
    Service->>Cache: write cache when implemented
  end
  Service-->>Handler: response DTO
  Handler-->>Client: standard JSON envelope
```

## Router

`server.SetupRoutes` creates `/api/v1` and registers feature routes.

## Middleware

Global middleware:

- Fiber recover middleware.
- CORS middleware with default config.
- Structured request logger.

Protected group middleware:

- JWT middleware.
- Permission middleware for admin routes.

## Validation

Handlers use `validator.ParseAndValidate[T]` for JSON body parsing and struct-tag validation.

## Controller

The codebase calls HTTP controllers `handlers`. Handlers are the only feature layer that imports Fiber.

## Service

Services hold business rules and interact with repositories, cache, storage URL generation, hashing, JWT generation, and worker jobs.

## Repository

Repositories use GORM with request contexts through `db.WithContext(ctx)`.

## Database

PostgreSQL is reached through GORM. Connection pooling is configured in `database.NewPostgresDB`.

## Response

All handlers return the shared response envelope from `internal/shared/response`.
