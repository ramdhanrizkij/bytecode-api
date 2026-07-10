# Security

## Authentication

JWT authentication is implemented with `middleware.JWTAuth`. Tokens must be sent as:

```text
Authorization: Bearer <token>
```

## Authorization

Permission-based RBAC is implemented with `middleware.RequirePermission`.

## Encryption

Transport encryption is not configured in application code. HTTPS termination is expected outside this codebase when required.

## Hashing

Passwords are hashed with bcrypt cost `12`.

Refresh tokens are hashed with SHA-256 before storage.

## CSRF

Not present in the analyzed codebase. The API uses Bearer tokens rather than cookie sessions.

## XSS

Not directly applicable to the JSON API. No HTML rendering layer is present.

## SQL Injection

Most queries use GORM parameter binding. The `sort` query parameter is interpolated into `Order` without a field allowlist, so callers can influence the order clause.

## Secrets

Secrets are loaded from environment variables. Production startup rejects the default JWT secret.

## HTTPS

Not configured in application code.

## Security Headers

Not present in the analyzed codebase.

## OWASP Considerations

| Area | Current State |
| --- | --- |
| Broken access control | RBAC middleware protects admin routes; `superadmin` bypass is hard-coded |
| Cryptographic failures | bcrypt for passwords, HMAC JWT, SHA-256 refresh token hashes |
| Injection | Parameter binding is used for filters; sort field needs allowlisting |
| Insecure design | Refresh-token rotation is implemented |
| Security misconfiguration | Swagger is protected in production only when credentials are configured; otherwise disabled |
| Identification failures | JWT expiration is checked by parser |
| Logging | Request and job logs exist; no audit log exists |

## Not Present

- Rate limiting.
- Account lockout.
- MFA.
- CSRF middleware.
- Security headers middleware.
- Audit trail table.
