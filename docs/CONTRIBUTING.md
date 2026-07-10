# Contributing

## Folder Conventions

Follow the existing feature-module layout:

```text
internal/features/<feature>/
  domain/
  handler/
  repository/
  service/
```

Shared infrastructure belongs under `internal/core`. Cross-feature helpers belong under `internal/shared` or `pkg` depending on whether they are application-specific.

## Naming Conventions

- HTTP handlers use `*HTTPHandler`.
- Constructors use `New...`.
- Domain contracts are named `...Service` and `...Repository`.
- Request DTOs end with `Request`.
- Response DTOs end with `Response`.
- Permission names use dotted names such as `users.view`.

## Architecture Rules

- Keep Fiber imports in handler and middleware layers.
- Keep GORM imports in repository and database infrastructure layers.
- Define service and repository contracts in feature `domain` packages.
- Use migrations for schema changes.
- Return standard response envelopes.
- Use `AppError` for expected application errors.

## Branch Naming

Not present in the analyzed codebase.

## Commit Convention

Not present in the analyzed codebase.

## PR Checklist

- Add or update migrations for schema changes.
- Add or update DTO validation tags for request changes.
- Update Swagger annotations for endpoint changes.
- Run `make test`.
- Run `make swagger` when API annotations change.
- Update docs when architecture, configuration, or API behavior changes.

## Review Checklist

- Verify handler, service, repository boundaries.
- Verify authorization middleware on protected routes.
- Verify cache invalidation for mutations.
- Verify errors use correct HTTP codes.
- Verify secrets are not committed.
- Verify migrations have up and down files.

## Definition of Done

- Code compiles.
- Tests pass or documented blockers exist.
- New endpoints have validation and Swagger annotations.
- New protected endpoints require JWT and permission middleware.
- Documentation reflects the implemented behavior.
