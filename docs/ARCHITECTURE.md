# Architecture

## System Context

```mermaid
flowchart LR
  Client[API Client] --> API[Bytecode API]
  API --> DB[(PostgreSQL)]
  API --> Redis[(Redis optional)]
  API --> Storage[Local Filesystem or MinIO]
  API --> Pool[In-process Worker Pool]
  Worker[Worker Process] --> DB
  Worker --> Scheduler[Scheduler]
  Scheduler --> Health[Health Check Job]
  Scheduler --> Cleanup[Refresh Token Cleanup]
  classDef frontend fill:#BBDEFB
  classDef backend fill:#C8E6C9
  classDef database fill:#FFE082
  classDef cache fill:#F8BBD0
  classDef queue fill:#D1C4E9
  classDef external fill:#CFD8DC
  class Client frontend
  class API,Worker,Scheduler,Health,Cleanup,Pool backend
  class DB database
  class Redis cache
  class Storage external
```

## Component Diagram

```mermaid
flowchart TD
  Main[cmd/api/main.go] --> Config[config.LoadConfig]
  Main --> Logger[pkg/logger]
  Main --> Database[database.NewPostgresDB]
  Main --> WorkerPool[worker.NewWorkerPool]
  Main --> Server[server.NewServer]
  Server --> Middleware[Global middleware]
  Server --> Auth[auth handler/service/repository]
  Server --> Users[user handler/service/repository]
  Server --> Roles[role handler/service/repository]
  Server --> Permissions[permission handler/service/repository]
  Server --> Cache[cache.Client]
  Server --> Storage[storage.Provider]
  Auth --> JWT[pkg/jwt]
  Auth --> Hash[pkg/hash]
  Users --> Hash
  Roles --> Cache
  Permissions --> Cache
  classDef frontend fill:#BBDEFB
  classDef backend fill:#C8E6C9
  classDef database fill:#FFE082
  classDef cache fill:#F8BBD0
  classDef queue fill:#D1C4E9
  classDef external fill:#CFD8DC
  class Main,Config,Logger,Server,Middleware,Auth,Users,Roles,Permissions,JWT,Hash backend
  class Database database
  class Cache cache
  class WorkerPool queue
  class Storage external
```

## Container Diagram

```mermaid
flowchart LR
  API[api container or process] --> Postgres[(postgres service)]
  API --> Redis[(redis service optional)]
  API --> MinIO[minio service optional]
  Worker[worker process] --> Postgres
  Compose[docker-compose.yml] --> Postgres
  Compose --> Redis
  Compose --> MinIO
  classDef frontend fill:#BBDEFB
  classDef backend fill:#C8E6C9
  classDef database fill:#FFE082
  classDef cache fill:#F8BBD0
  classDef queue fill:#D1C4E9
  classDef external fill:#CFD8DC
  class API,Worker,Compose backend
  class Postgres database
  class Redis cache
  class MinIO external
```

## Module Diagram

```mermaid
flowchart TD
  Handler[handler] --> Service[service]
  Service --> Repository[repository]
  Repository --> Model[internal/model]
  Handler --> Shared[shared response validator pagination errors]
  Service --> Core[core cache storage worker]
  Repository --> GORM[GORM DB]
  classDef frontend fill:#BBDEFB
  classDef backend fill:#C8E6C9
  classDef database fill:#FFE082
  classDef cache fill:#F8BBD0
  classDef queue fill:#D1C4E9
  classDef external fill:#CFD8DC
  class Handler,Service,Repository,Model,Shared,Core backend
  class GORM database
```

## Dependency Graph

- `cmd/api` depends on `internal/core/*`, generated `docs`, and `pkg/logger`.
- `internal/core/server` manually constructs repositories, services, handlers, cache, and storage.
- Feature handlers import Fiber and feature domain interfaces.
- Feature services import domain interfaces, shared errors, and infrastructure interfaces where needed.
- Feature repositories import GORM and `internal/model`.
- Domain packages define DTOs and interfaces. They do not import Fiber.

## Design Patterns

| Pattern | Evidence |
| --- | --- |
| Feature modules | `internal/features/auth`, `user`, `role`, `permission` |
| Repository pattern | `domain.*Repository` interfaces and GORM implementations |
| Service layer | `domain.*Service` interfaces implemented under `service/` |
| Manual dependency injection | `server.SetupRoutes` constructs dependencies explicitly |
| Standard response envelope | `internal/shared/response` |
| Middleware pipeline | Fiber middleware for recovery, CORS, logging, JWT, RBAC |
| Worker pool | `internal/core/worker/pool.go` |
| Scheduler | `internal/core/worker/scheduler.go` |

## Architectural Decisions

- Database schema is migration-first. `AutoMigrate` is intentionally not used.
- The API and scheduled worker are separate executables.
- JWT access tokens carry `user_id`, `email`, and `role_name`.
- Refresh tokens are opaque random values stored only as SHA-256 hashes.
- Redis is optional; disabled mode uses a no-op cache client.
- Storage is selected at runtime by `STORAGE_PROVIDER`.

## Communication Protocols

- HTTP JSON API under `/api/v1`.
- PostgreSQL wire protocol through GORM.
- Redis protocol through `redis/go-redis` when enabled.
- MinIO S3-compatible API when `STORAGE_PROVIDER=minio`.

## Synchronization Strategy

- Worker pool uses a buffered Go channel and goroutines.
- Scheduler uses `time.Ticker`, context cancellation, and `sync.WaitGroup`.
- RBAC middleware uses an in-memory `sync.Map` for role permission cache.
- Service caches are invalidated by Redis key prefixes.

## Data Ownership

| Data | Owner |
| --- | --- |
| Users | User and Auth features |
| Roles | Role feature |
| Permissions | Permission feature |
| Role assignments | Role feature through `role_permissions` |
| Refresh tokens | Auth feature |
| Profile picture references | User feature |

## Cross-Cutting Concerns

- Logging: Zap request logs and job logs.
- Error handling: `AppError` plus Fiber `ErrorHandler`.
- Validation: `go-playground/validator`.
- Security: JWT authentication, RBAC authorization, bcrypt hashing.
- Caching: Redis/no-op abstraction and middleware in-memory cache.
