# Testing

## Unit Tests

Unit tests live under `test/unit`.

Implemented coverage:

- Role service cache read behavior.
- Permission service cache invalidation.
- User permission cache behavior.
- User cache invalidation on delete.
- Local storage put/delete.
- Storage bucket configuration.
- Auth service tests.
- User service tests.

## Integration Tests

Integration tests live under `test/integration` and use Testcontainers PostgreSQL.

Implemented tests:

- `TestHealthCheck`
- `TestAuth_Register`
- `TestAuth_RefreshAndLogout`
- `TestRole_CRUD`
- `TestUser_ProfileLifecycle`

## E2E

Browser or external E2E tests are not present in the analyzed codebase.

## Coverage

No coverage threshold or coverage report command is configured in the Makefile.

## Fixtures

Integration tests create a PostgreSQL container, run migrations, and truncate tables between cases through `test/integration/testhelper`.

## Mocking

Unit tests use Go test patterns and Testify. Detailed mocking strategy depends on individual test files.

## Test Commands

```bash
make test
make test-unit
make test-integration
```

`make test-integration` runs:

```bash
go test -v -tags=integration ./test/integration/...
```

## External Requirements

Integration tests require Docker because Testcontainers starts PostgreSQL.
