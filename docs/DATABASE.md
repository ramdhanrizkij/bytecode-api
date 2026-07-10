# Database

## Database Engine

PostgreSQL is the only implemented database engine. The application connects through GORM using `gorm.io/driver/postgres`.

## Schema

| Table | Purpose |
| --- | --- |
| `roles` | Role definitions such as `superadmin`, `admin`, and `user` |
| `permissions` | Permission names such as `users.view` |
| `users` | User accounts and role assignment |
| `role_permissions` | Many-to-many role permission join table |
| `refresh_tokens` | Hashed opaque refresh tokens |

## Relationships

```mermaid
erDiagram
  roles ||--o{ users : assigned_to
  roles ||--o{ role_permissions : has
  permissions ||--o{ role_permissions : grants
  users ||--o{ refresh_tokens : owns
  roles {
    uuid id PK
    varchar name UK
    text description
    varchar guard_name
    timestamptz created_at
    timestamptz updated_at
  }
  permissions {
    uuid id PK
    varchar name UK
    text description
    varchar guard_name
    timestamptz created_at
    timestamptz updated_at
  }
  users {
    uuid id PK
    varchar name
    varchar email UK
    varchar password
    text profile_picture
    uuid role_id FK
    boolean is_active
    timestamptz created_at
    timestamptz updated_at
  }
  role_permissions {
    uuid id PK
    uuid role_id FK
    uuid permission_id FK
    timestamptz created_at
  }
  refresh_tokens {
    uuid id PK
    uuid user_id FK
    varchar token_hash UK
    timestamptz expires_at
    timestamptz revoked_at
    timestamptz created_at
    timestamptz updated_at
  }
```

## Indexes

| Table | Index |
| --- | --- |
| `roles` | unique `name` |
| `permissions` | unique `name` |
| `users` | unique `email` |
| `role_permissions` | unique `(role_id, permission_id)` |
| `refresh_tokens` | unique `token_hash` |
| `refresh_tokens` | `idx_refresh_tokens_user_id` |
| `refresh_tokens` | `idx_refresh_tokens_expires_at` |

## Constraints

- Primary keys use UUID with `gen_random_uuid()`.
- `users.role_id` references `roles(id)` with `ON DELETE SET NULL`.
- `role_permissions.role_id` references `roles(id)` with `ON DELETE CASCADE`.
- `role_permissions.permission_id` references `permissions(id)` with `ON DELETE CASCADE`.
- `refresh_tokens.user_id` references `users(id)` with `ON DELETE CASCADE`.

## Migration Strategy

Migrations live under `migrations/` and are applied with `golang-migrate`.

```bash
make migrate-up
make migrate-down
make migrate-refresh
make migrate-create name=create_something_table
```

GORM `AutoMigrate` is not used.

## Seed Data

Migration `000005` inserts roles:

- `superadmin`
- `admin`
- `user`

Migration `000006` inserts permissions and assigns all permissions to `superadmin`.

## Soft Deletes

Not present in the analyzed codebase. Models do not include `gorm.DeletedAt`, and deletes execute hard deletes.

## Transactions

Explicit multi-step transaction blocks are not present in the analyzed codebase.

## Isolation

No custom transaction isolation level is configured.
