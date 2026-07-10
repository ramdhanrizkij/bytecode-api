# Validation

## Validation Layer

Validation lives in `internal/shared/validator`.

## Libraries

The code uses `github.com/go-playground/validator/v10`.

## Request Validation

Handlers call:

```go
validator.ParseAndValidate[RequestType](c)
```

This function:

1. Parses the JSON body using Fiber binding.
2. Validates struct tags.
3. Writes `400` on parse failure.
4. Writes `422` on validation failure.
5. Returns the typed request on success.

## Business Validation

Services enforce business rules:

- Duplicate email checks.
- Duplicate role name checks.
- Duplicate permission name checks.
- Invalid UUID checks after parsing string fields.
- Active user check on login.
- Self-delete prevention.
- Refresh-token revocation and expiration checks.

## Sanitization

Profile picture strings are trimmed and empty strings become `nil`.

Storage object keys normalize slashes, strip leading slashes, and remove empty, `.`, and `..` path segments.

## Normalization

- Pagination defaults are normalized in `pagination.NewPaginationQuery`.
- Empty `guard_name` defaults to `api` on create for roles and permissions.
- Empty storage provider defaults to `local`.

## Not Present

- Request body size validation per endpoint. Fiber has a global `4MB` body limit.
- HTML sanitization.
- SQL identifier allowlisting for the `sort` query parameter.
