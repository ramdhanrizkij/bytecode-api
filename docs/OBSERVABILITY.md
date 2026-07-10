# Observability

## Logging

Zap logging is initialized in `pkg/logger`.

| Level | Encoder |
| --- | --- |
| `debug` | Development console encoder |
| `info`, `warn`, `error` | Production JSON encoder |

Logs write to stdout, with errors to stderr.

## Request Logs

`middleware.RequestLogger` logs:

- HTTP method.
- Path.
- Response status.
- Latency.
- Client IP.

## Metrics

Not present in the analyzed codebase.

## Tracing

Not present in the analyzed codebase.

## Health Checks

HTTP health endpoint:

```text
GET /api/v1/health
```

It reports service name, environment, database status, cache status, and storage provider.

Worker health job:

- Name: `system_health_check`.
- Interval: 5 minutes.
- Behavior: pings database and logs result.

## Monitoring

Not present in the analyzed codebase.

## Alerting

Not present in the analyzed codebase.
