# Authentication

## Authentication Flow

The API uses email/password authentication, JWT access tokens, and opaque refresh tokens.

```mermaid
sequenceDiagram
  participant Client
  participant AuthHandler
  participant AuthService
  participant AuthRepository
  participant DB as PostgreSQL
  participant WorkerPool
  Client->>AuthHandler: POST /api/v1/auth/login
  AuthHandler->>AuthService: Login(email, password)
  AuthService->>AuthRepository: FindUserByEmail(email)
  AuthRepository->>DB: SELECT users + role
  DB-->>AuthRepository: User
  AuthService->>AuthService: bcrypt password check
  AuthService->>AuthService: Generate JWT
  AuthService->>AuthRepository: CreateRefreshToken(hash)
  AuthRepository->>DB: INSERT refresh_tokens
  AuthService-->>AuthHandler: AuthResponse
  AuthHandler-->>Client: token + refresh_token
```

## JWT

Access tokens are generated in `pkg/jwt`.

| Attribute | Value |
| --- | --- |
| Signing method | HMAC SHA-256 |
| Secret | `JWT_SECRET` |
| Claims | `user_id`, `email`, `role_name`, `iat`, `exp` |
| Expiry | `JWT_EXPIRY_HOURS`, default `24` |

The parser rejects unexpected signing methods.

## OAuth

Not present in the analyzed codebase.

## Session

Server-side web sessions are not present in the analyzed codebase.

## Refresh Token

Refresh tokens are random 32-byte values encoded as hex. Only SHA-256 hashes are stored in `refresh_tokens.token_hash`.

## Token Lifetime

| Token | Configuration | Default |
| --- | --- | --- |
| Access token | `JWT_EXPIRY_HOURS` | `24` hours |
| Refresh token | `JWT_REFRESH_EXPIRY_HOURS` | `168` hours |

## Middleware

`middleware.JWTAuth` reads the `Authorization` header, requires `Bearer <token>`, parses the token, and stores claims in `c.Locals("user")`.

## Login Sequence

```mermaid
sequenceDiagram
  participant Client
  participant API
  participant DB
  Client->>API: POST /auth/login
  API->>DB: Find user by email with role
  DB-->>API: User
  API->>API: bcrypt.CheckPassword
  API->>API: reject if inactive
  API->>API: sign access JWT
  API->>DB: insert hashed refresh token
  API-->>Client: access token and refresh token
```

## Logout Sequence

```mermaid
sequenceDiagram
  participant Client
  participant API
  participant DB
  Client->>API: POST /auth/logout
  API->>API: SHA-256 hash refresh token
  API->>DB: Find refresh token by hash
  alt found
    API->>DB: set revoked_at
  else not found
    API-->>Client: 200 logout successful
  end
  API-->>Client: 200 logout successful
```
