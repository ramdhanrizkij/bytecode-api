# Error Handling

## Error Type

`internal/shared/errors.AppError` is the standard application error. It stores:

- HTTP status code.
- User-facing message.
- Wrapped internal cause.

## Sentinel Errors

| Error | Code | Message |
| --- | --- | --- |
| `ErrUnauthorized` | 401 | `unauthorized` |
| `ErrForbidden` | 403 | `forbidden` |
| `ErrNotFound` | 404 | `resource not found` |
| `ErrConflict` | 409 | `conflict` |
| `ErrValidation` | 422 | `validation error` |
| `ErrInternalServer` | 500 | `internal server error` |

## HTTP Mapping

Handlers use `errors.As` to detect `*AppError` and call `response.Error`.

## Global Error Handler

Fiber is configured with `customErrorHandler` in `server.NewServer`. It converts `AppError` into the standard response envelope. Unknown errors return `500`.

## Panic Recovery

The server uses Fiber's built-in `recover.New()`. A custom `middleware.Recovery` also exists, but it is not registered by `server.NewServer`.

## Validation Errors

Validation failures return:

```json
{
  "meta": {
    "code": 422,
    "message": "validation error"
  },
  "errors": [
    {
      "field": "Email",
      "tag": "email",
      "message": "Format email tidak valid"
    }
  ]
}
```

## Not Present

- Centralized error codes beyond HTTP status code.
- Problem Details RFC 7807 format.
- Error correlation IDs.
