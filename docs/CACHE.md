# Cache

## Redis

Redis is optional. `cache.NewClient` returns:

- Redis-backed client when `REDIS_ENABLED=true`.
- No-op disabled client when Redis is disabled.

## Memory Cache

RBAC middleware has an in-memory `sync.Map` cache for role permission names. The TTL is `5 minutes`.

## TTL

Service cache TTL comes from `REDIS_CACHE_TTL_MINUTES`, default `5`.

## Cache Keys

| Feature | Key Pattern |
| --- | --- |
| Users list | `users:list:page=<n>:per_page=<n>:sort=<field>:order=<dir>:search=<term>` |
| User detail | `users:detail:<id>` |
| User permissions | `users:permissions:<id>` |
| Roles list | `roles:list:page=<n>:per_page=<n>:sort=<field>:order=<dir>:search=<term>` |
| Role detail | `roles:detail:<id>` |
| Permissions list | `permissions:list:page=<n>:per_page=<n>:sort=<field>:order=<dir>:search=<term>` |
| Permission detail | `permissions:detail:<id>` |

## Cache Invalidation

| Mutation | Invalidation |
| --- | --- |
| User create/update/delete | Delete Redis keys with prefix `users:` |
| Role create/update/delete/permission assignment/removal | Delete Redis keys with prefix `roles:` |
| Permission create | Delete Redis keys with prefix `permissions:` |
| Permission update/delete | Delete Redis keys with prefixes `permissions:` and `roles:` |

## Patterns

- Cache-aside read path.
- JSON marshal/unmarshal for cached values.
- Redis `SCAN` plus `DEL` for prefix invalidation.
- No-op cache interface when disabled.

## Not Present

- Distributed invalidation for the RBAC middleware in-memory cache.
- Cache metrics.
- Cache stampede protection.
