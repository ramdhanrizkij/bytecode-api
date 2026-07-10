# Authorization

## RBAC

Role-based access control is implemented with:

- `roles`
- `permissions`
- `role_permissions`
- `middleware.RequirePermission`
- JWT `role_name` claim

The middleware reads the authenticated user's role name from JWT claims, loads permissions for that role, and checks whether at least one required permission exists.

## ABAC

Not present in the analyzed codebase.

## Policies

No standalone policy objects or policy engine are present. Authorization rules are attached directly in route registration files.

## Permissions

Seeded permission names:

- `roles.view`
- `roles.create`
- `roles.edit`
- `roles.delete`
- `roles.assign-permission`
- `roles.remove-permission`
- `permissions.view`
- `permissions.create`
- `permissions.edit`
- `permissions.delete`
- `users.view`
- `users.create`
- `users.edit`
- `users.delete`

## Guards

`guard_name` exists on roles and permissions, defaulting to `api`. The analyzed authorization middleware does not filter by `guard_name`.

## Middleware

| Middleware | Responsibility |
| --- | --- |
| `JWTAuth` | Validates Bearer JWT and stores claims |
| `RequirePermission` | Checks permission names for the user's role |
| `RequireRole` | Checks role names; defined but not used by current routes |

## Role Hierarchy

There is no persisted hierarchy. `superadmin` has a hard-coded bypass in `RequirePermission`.

## Permission Resolution

```mermaid
flowchart TD
  Request[Protected request] --> JWT[JWTAuth]
  JWT --> Claims[Claims role_name]
  Claims --> Superadmin{role_name is superadmin}
  Superadmin -->|yes| Allow[Allow request]
  Superadmin -->|no| Cache[Check in-memory role permission cache]
  Cache --> DB[Load Role with Permissions if cache miss]
  DB --> Compare[Compare required permission names]
  Compare -->|match| Allow
  Compare -->|no match| Deny[403 insufficient permissions]
  classDef frontend fill:#BBDEFB
  classDef backend fill:#C8E6C9
  classDef database fill:#FFE082
  classDef cache fill:#F8BBD0
  classDef queue fill:#D1C4E9
  classDef external fill:#CFD8DC
  class Request frontend
  class JWT,Claims,Superadmin,Compare,Allow,Deny backend
  class Cache cache
  class DB database
```
