# Events

## Domain Events

Not present in the analyzed codebase.

## Application Events

No event bus abstraction exists. The closest implemented asynchronous action is direct submission of a welcome-email job after user registration.

## Event Bus

Not present in the analyzed codebase.

## Subscribers

Not present in the analyzed codebase.

## Publishers

Not present as event publishers. `authService.Register` directly calls `workerPool.Submit`.

## Flow Diagram

```mermaid
flowchart TD
  Register[authService.Register] --> CreateUser[Create user]
  CreateUser --> SubmitJob[Submit SendWelcomeEmailJob]
  SubmitJob --> WorkerPool[Worker pool channel]
  WorkerPool --> Execute[Job Execute]
  Execute --> Log[Log simulated email sent]
  classDef frontend fill:#BBDEFB
  classDef backend fill:#C8E6C9
  classDef database fill:#FFE082
  classDef cache fill:#F8BBD0
  classDef queue fill:#D1C4E9
  classDef external fill:#CFD8DC
  class Register,CreateUser,Execute,Log backend
  class SubmitJob,WorkerPool queue
```

## Messaging

External messaging systems such as Kafka, RabbitMQ, NATS, or SQS are not present in the analyzed codebase.
