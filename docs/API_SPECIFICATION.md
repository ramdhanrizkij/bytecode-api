# API Specification

## Base URL

All application endpoints are mounted under:

```text
/api/v1
```

Swagger UI is mounted separately at `/swagger/*` when enabled.

## Response Envelope

```json
{
  "meta": {
    "code": 200,
    "message": "Message"
  },
  "data": {}
}
```

Paginated responses add:

```json
{
  "pagination": {
    "current_page": 1,
    "per_page": 10,
    "total_items": 42,
    "total_pages": 5
  }
}
```

## Common Headers

| Header | Required | Applies To |
| --- | --- | --- |
| `Content-Type: application/json` | Yes for JSON bodies | POST, PUT, DELETE with body |
| `Authorization: Bearer <token>` | Yes for protected endpoints | Users, Roles, Permissions |

## Pagination, Filtering, Sorting

List endpoints for users, roles, and permissions support:

| Query | Default | Rules |
| --- | --- | --- |
| `page` | `1` | Values below `1` become `1` |
| `per_page` | `10` | Maximum `100` |
| `sort` | `created_at` | Passed to GORM as a column name |
| `order` | `desc` | Only `asc` or `desc`; other values become `desc` |
| `search` | empty | Uses PostgreSQL `ILIKE` on implemented fields |

Rate limiting: Not present in the analyzed codebase.

## GET /health

### Description

Returns API, database, cache, and storage health status.

### Authentication

None.

### Response

```json
{
  "meta": {
    "code": 200,
    "message": "service is healthy"
  },
  "data": {
    "status": "ok",
    "service": "bytecode-api",
    "environment": "development",
    "database": "up",
    "cache": "disabled",
    "storage": "local"
  }
}
```

### Status Codes

| Code | Meaning |
| --- | --- |
| 200 | Health data returned. `data.status` may be `degraded` if database ping fails. |

## POST /auth/register

### Description

Creates a user with the default `user` role, submits a welcome-email job, and returns an access token plus refresh token.

### Authentication

None.

### Request Body

```json
{
  "name": "Jane Doe",
  "email": "jane@example.com",
  "password": "secret123",
  "profile_picture": "profiles/jane.png"
}
```

### Validation Rules

| Field | Rules |
| --- | --- |
| `name` | required, min 2, max 100 |
| `email` | required, email |
| `password` | required, min 8, max 50 |
| `profile_picture` | optional, max 500 |

### Response

```json
{
  "meta": {
    "code": 201,
    "message": "User registered successfully"
  },
  "data": {
    "user": {
      "id": "018f7606-a3f7-7c40-8e4b-2d47c6e04c8d",
      "name": "Jane Doe",
      "email": "jane@example.com",
      "role_name": "user",
      "is_active": true
    },
    "token": "access-token-value",
    "refresh_token": "refresh-token-value"
  }
}
```

### Possible Errors

| Code | Reason |
| --- | --- |
| 400 | Malformed JSON |
| 409 | Email already registered |
| 422 | Validation failure |
| 500 | Default role missing, database error, hashing error, token generation error |

## POST /auth/login

### Description

Authenticates an active user by email and password.

### Authentication

None.

### Request Body

```json
{
  "email": "jane@example.com",
  "password": "secret123"
}
```

### Validation Rules

| Field | Rules |
| --- | --- |
| `email` | required, email |
| `password` | required |

### Response

Same shape as `/auth/register` with status `200`.

### Possible Errors

| Code | Reason |
| --- | --- |
| 400 | Malformed JSON |
| 401 | User not found or password mismatch |
| 403 | User exists but `is_active` is false |
| 422 | Validation failure |
| 500 | Database or token generation error |

## POST /auth/refresh

### Description

Exchanges a valid refresh token for a new access token and a rotated refresh token. The old refresh token is revoked.

### Authentication

None.

### Request Body

```json
{
  "refresh_token": "refresh-token-value"
}
```

### Validation Rules

| Field | Rules |
| --- | --- |
| `refresh_token` | required |

### Response

```json
{
  "meta": {
    "code": 200,
    "message": "Token refreshed successfully"
  },
  "data": {
    "token": "new-access-token",
    "refresh_token": "new-refresh-token"
  }
}
```

### Possible Errors

| Code | Reason |
| --- | --- |
| 401 | Token hash not found, revoked token, expired token, inactive user |
| 422 | Validation failure |
| 500 | Database or token generation error |

## POST /auth/logout

### Description

Revokes a refresh token. Missing token records are treated as successful logout.

### Authentication

None.

### Request Body

```json
{
  "refresh_token": "refresh-token-value"
}
```

### Response

```json
{
  "meta": {
    "code": 200,
    "message": "Logout successful"
  }
}
```

### Possible Errors

| Code | Reason |
| --- | --- |
| 422 | Validation failure |
| 500 | Database error |

## GET /users/me

### Description

Returns the authenticated user's profile.

### Authentication

JWT required.

### Required Permission

None beyond a valid JWT.

### Response

```json
{
  "meta": {
    "code": 200,
    "message": "Current user profile retrieved successfully"
  },
  "data": {
    "id": "018f7606-a3f7-7c40-8e4b-2d47c6e04c8d",
    "name": "Jane Doe",
    "email": "jane@example.com",
    "is_active": true,
    "role": {
      "id": "018f7606-a3f7-7c40-8e4b-2d47c6e04c8d",
      "name": "user"
    },
    "created_at": "2026-05-13T10:00:00+07:00",
    "updated_at": "2026-05-13T10:00:00+07:00"
  }
}
```

### Possible Errors

| Code | Reason |
| --- | --- |
| 401 | Missing or invalid JWT |
| 404 | User not found |
| 500 | Database error |

## PUT /users/me/profile

### Description

Updates the authenticated user's name, email, and profile-picture reference.

### Authentication

JWT required.

### Request Body

```json
{
  "name": "Jane Doe",
  "email": "jane@example.com",
  "profile_picture": "profiles/jane.png"
}
```

### Validation Rules

| Field | Rules |
| --- | --- |
| `name` | required, min 2, max 100 |
| `email` | required, email |
| `profile_picture` | optional, max 500 |

### Possible Errors

| Code | Reason |
| --- | --- |
| 400 | Malformed JSON |
| 401 | Missing or invalid JWT |
| 409 | Email already exists |
| 422 | Validation failure |
| 500 | Database error |

## GET /users/me/permissions

### Description

Returns permission names assigned to the authenticated user's role.

### Authentication

JWT required.

### Response

```json
{
  "meta": {
    "code": 200,
    "message": "User permissions retrieved successfully"
  },
  "data": [
    "users.view"
  ]
}
```

## GET /users/

### Description

Returns paginated users. Search matches `name` or `email` with `ILIKE`.

### Authentication

JWT required.

### Required Permission

`users.view`, except `superadmin` bypasses permission checks.

### Possible Errors

| Code | Reason |
| --- | --- |
| 401 | Missing or invalid JWT |
| 403 | Missing `users.view` |
| 500 | Database error |

## GET /users/{id}

### Description

Returns one user by UUID.

### Authentication

JWT required.

### Required Permission

`users.view`.

## POST /users/

### Description

Creates a user with a specified role.

### Authentication

JWT required.

### Required Permission

`users.create`.

### Request Body

```json
{
  "name": "Jane Doe",
  "email": "jane@example.com",
  "password": "secret123",
  "role_id": "018f7606-a3f7-7c40-8e4b-2d47c6e04c8d",
  "profile_picture": "profiles/jane.png",
  "is_active": true
}
```

### Validation Rules

| Field | Rules |
| --- | --- |
| `name` | required, min 2, max 100 |
| `email` | required, email |
| `password` | required, min 8, max 50 |
| `role_id` | required, uuid |
| `profile_picture` | optional, max 500 |
| `is_active` | optional boolean |

## PUT /users/{id}

### Description

Updates a user. Password is optional; when provided it is re-hashed with bcrypt.

### Required Permission

`users.edit`.

### Request Body

Same as create user, except `password` is optional with `omitempty,min=8,max=50`.

## DELETE /users/{id}

### Description

Deletes a user. A user cannot delete their own account.

### Required Permission

`users.delete`.

### Possible Errors

| Code | Reason |
| --- | --- |
| 403 | Missing permission or deleting own account |
| 404 | Target user not found |

## GET /roles/

### Description

Returns paginated roles. Search matches `name` with `ILIKE`.

### Required Permission

`roles.view`.

## GET /roles/{id}

### Description

Returns one role by UUID with permissions preloaded.

### Required Permission

`roles.view`.

## POST /roles/

### Description

Creates a role. Empty `guard_name` defaults to `api`.

### Required Permission

`roles.create`.

### Request Body

```json
{
  "name": "manager",
  "description": "Can manage operational resources",
  "guard_name": "api"
}
```

### Validation Rules

| Field | Rules |
| --- | --- |
| `name` | required, min 2, max 50 |
| `description` | max 255 |
| `guard_name` | max 50 |

## PUT /roles/{id}

### Description

Updates a role. Empty `guard_name` leaves the existing value unchanged.

### Required Permission

`roles.edit`.

## DELETE /roles/{id}

### Description

Deletes a role. Database foreign keys set user `role_id` to null and cascade role-permission rows.

### Required Permission

`roles.delete`.

## POST /roles/{id}/permissions

### Description

Assigns permissions to a role. Duplicate role-permission pairs are ignored with `ON CONFLICT DO NOTHING`.

### Required Permission

`roles.assign-permission`.

### Request Body

```json
{
  "permission_ids": [
    "018f7606-a3f7-7c40-8e4b-2d47c6e04c8d"
  ]
}
```

### Validation Rules

| Field | Rules |
| --- | --- |
| `permission_ids` | required, min 1, each item uuid |

## DELETE /roles/{id}/permissions

### Description

Removes permission assignments from a role.

### Required Permission

`roles.remove-permission`.

### Request Body

Same as `POST /roles/{id}/permissions`.

## GET /permissions/

### Description

Returns paginated permissions. Search matches `name` with `ILIKE`.

### Required Permission

`permissions.view`.

## GET /permissions/{id}

### Description

Returns one permission by UUID.

### Required Permission

`permissions.view`.

## POST /permissions/

### Description

Creates a permission. Empty `guard_name` defaults to `api`.

### Required Permission

`permissions.create`.

### Request Body

```json
{
  "name": "users.view",
  "description": "Can view users",
  "guard_name": "api"
}
```

### Validation Rules

| Field | Rules |
| --- | --- |
| `name` | required, min 2, max 100 |
| `description` | max 255 |
| `guard_name` | max 50 |

## PUT /permissions/{id}

### Description

Updates a permission. Permission cache is invalidated, and role cache is also invalidated.

### Required Permission

`permissions.edit`.

## DELETE /permissions/{id}

### Description

Deletes a permission. Role-permission rows are removed by cascade.

### Required Permission

`permissions.delete`.
