# Domain Model

## Entities

| Entity | Source | Description |
| --- | --- | --- |
| User | `internal/model/user.go` | Account with email, password hash, active flag, optional role, and optional profile picture key |
| Role | `internal/model/role.go` | Named role with guard name and permissions |
| Permission | `internal/model/permission.go` | Named capability such as `users.view` |
| RolePermission | `internal/model/role_permission.go` | Explicit join entity between roles and permissions |
| RefreshToken | `internal/model/refresh_token.go` | Hashed refresh token linked to a user |

## Value Objects

Not present as explicit value-object types. UUIDs, email addresses, permission names, and profile-picture keys are represented as primitive strings or `uuid.UUID`.

## Aggregates

No aggregate pattern is explicitly named. Practical consistency boundaries are:

- User with Role projection for profile and list responses.
- Role with Permissions for RBAC and role detail responses.
- RefreshToken with User and User.Role for token refresh.

## Repositories

Feature domain packages define repository interfaces:

- `AuthRepository`
- `UserRepository`
- `RoleRepository`
- `PermissionRepository`

GORM implementations live under each feature's `repository/` directory.

## Factories

Dedicated factory objects are not present. Constructors exist for handlers, services, repositories, cache, storage, workers, and logger.

## Domain Services

Feature service implementations contain business rules:

- Auth: registration, login, refresh-token rotation, logout, token cleanup.
- User: duplicate email checks, password hashing, cache invalidation, self-delete prevention.
- Role: duplicate name checks, permission assignment/removal.
- Permission: duplicate name checks, cache invalidation.

## Bounded Contexts

The codebase uses feature folders rather than formal DDD bounded context packages:

- Auth
- User
- Role
- Permission

## Business Rules

- Registration assigns the default `user` role by name.
- `superadmin` bypasses permission checks.
- Users cannot delete their own account.
- Inactive users cannot log in.
- Refresh tokens are single-use after rotation.
- Empty role and permission `guard_name` values default to `api` on create.
