# Glossary

## Access Token

JWT returned by login, registration, or refresh. It authenticates protected API requests.

## AppError

Application error type with HTTP status code, public message, and wrapped cause.

## Guard Name

String field on roles and permissions. Defaults to `api`. It is stored but not used by authorization checks.

## Permission

Named capability such as `users.view` or `roles.create`.

## RBAC

Role-based access control. Users have roles. Roles have permissions. Routes require permission names.

## Refresh Token

Opaque random token used to issue a new access token. Only its SHA-256 hash is stored.

## Role

Named group of permissions, such as `superadmin`, `admin`, or `user`.

## Superadmin

Role name with hard-coded authorization bypass in `RequirePermission`.

## Worker Pool

In-process queue and goroutine pool used by the API process for immediate asynchronous jobs.

## Scheduler

Worker-process component that runs registered tasks at fixed intervals.

## Storage Provider

Runtime abstraction for local filesystem or MinIO object storage.
