# Queues

## Jobs

| Job | Location | Execution Mode |
| --- | --- | --- |
| `SendWelcomeEmailJob` | `internal/features/auth/job/send_email_job.go` | API process worker pool |
| `system_health_check` | `internal/core/worker/jobs/health_check.go` | Worker process scheduler every 5 minutes |
| `cleanup_expired_tokens` | `cmd/worker/main.go` through `AuthService.CleanupExpiredTokens` | Worker process scheduler every 1 hour |
| `CleanupExpiredTokensJob` | `internal/features/auth/job/cleanup_token_job.go` | Defined but not registered directly |

## Workers

`WorkerPool` starts a fixed number of goroutines. The API process configures `5` workers and queue size `100`.

## Retries

Not present in the analyzed codebase.

## Dead Letter Queue

Not present in the analyzed codebase.

## Scheduling

The worker process uses `Scheduler` with `time.Ticker`.

```mermaid
flowchart TD
  WorkerMain[cmd/worker/main.go] --> Scheduler[NewScheduler]
  Scheduler --> Health[Register health check every 5m]
  Scheduler --> Cleanup[Register token cleanup every 1h]
  Scheduler --> Ticker[time.Ticker per task]
  Ticker --> Execute[Execute task]
  Execute --> Log[Log success or error]
  classDef frontend fill:#BBDEFB
  classDef backend fill:#C8E6C9
  classDef database fill:#FFE082
  classDef cache fill:#F8BBD0
  classDef queue fill:#D1C4E9
  classDef external fill:#CFD8DC
  class WorkerMain,Scheduler,Health,Cleanup,Ticker,Execute,Log backend
```

## Concurrency

- Worker pool concurrency equals configured worker count.
- Scheduler runs each registered task in its own goroutine.
- Active jobs finish during worker-pool shutdown after the queue is closed.
