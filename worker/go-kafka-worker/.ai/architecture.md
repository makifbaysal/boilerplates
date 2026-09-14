# Architecture

Hexagonal (ports & adapters), same shape as the `backend/go-fiber`
boilerplate: the core knows nothing about Kafka, RabbitMQ, or any
specific database.

```
cmd/
  kafka/main.go             wiring: config -> logger -> store -> service -> kafka consumer -> Run
  rabbitmq/main.go            same wiring, mounts the RabbitMQ consumer instead
internal/
  core/job/                   the hexagon: no franz-go/amqp091-go/db-driver imports, ever
    job.go                     domain struct
    ports.go                   Store (outbound port) + Processor (inbound port) + sentinel errors
    service.go                 Processor impl — the only thing that calls Store
    mocks/                     mockery output for Store, generated — never edit by hand
  adapter/
    inbound/                  driving adapters: decoded message -> Processor call
      kafka/                   franz-go — see .ai/common-tasks.md
      rabbitmq/                 amqp091-go
    outbound/                 driven adapters: Store port -> real storage
      memory/                  default for cmd/{kafka,rabbitmq}, and what adapter tests run against
      postgres/, mongo/, cassandra/   build on pkg/{postgres,mongo,cassandra}
  platform/
    logger/                   log/slog (stdlib) + WithLogging message-handler middleware
    config/                    koanf loader, env-based demo configs + live reload
pkg/
  postgres/, mongo/, cassandra/   reusable, app-agnostic client constructors — copied
                                   verbatim from backend/go-fiber/pkg, zero task-specific code
  redis/, elasticsearch/          same idea, not wired to a Store (see backend/go-fiber's note)
configs/
  config.dev.yaml, config.staging.yaml, config.prod.yaml
```

## Message contract

Both consumers expect a JSON body: `{"task_id": "...", "type": "..."}`.
`payload` passed to `Processor.Process` is the raw message bytes, stored
verbatim by `Store.Record` — extend the `message` struct in
`adapter/inbound/{kafka,rabbitmq}/consumer.go` together if the contract
grows a field, and keep both in sync since they consume the same topic
shape from two different brokers.

## Layering rule

Same as go-fiber: `adapter/inbound/*` never touches `Store` directly —
always through `job.Processor`. `core/job` never imports a framework
package. Each consumer's `HandleMessage(ctx, source, value []byte) error`
is the actual per-message logic and is directly unit-testable without a
broker; `Run` is the thin polling/channel loop around it, covered by a
skip-gated integration test instead (see `.ai/testing.md`).

## Delivery semantics (deliberately simple)

Both consumers Ack/commit even when `HandleMessage` fails — a bad or
unprocessable message is logged and skipped, not retried or
dead-lettered. That's a boilerplate default, not a recommendation: wire
retry/DLQ policy into `Run` (or wrap `HandleMessage`) for a real
at-least-once use case.

## Logging

Same pattern as go-fiber's REST middleware / gRPC interceptor, adapted
for message consumers: `logger.WithLogging` wraps `HandleMessage` with
one structured log line per message (source, size, duration, error) and
attaches a request-scoped logger (`message_id`) to `ctx`, retrievable via
`logger.FromContext`.

## Config

Identical mechanism to go-fiber's `platform/config` (koanf,
`configs/config.<APP_ENV>.yaml` + `APP_`-prefixed env overrides + live
file-watch reload wired into a `slog.LevelVar`) — see
`backend/go-fiber/.ai/architecture.md#config` for the full write-up,
including why the Docker images don't default `APP_ENV`. Fields differ:
`Kafka{Brokers,Topic,GroupID}` and `RabbitMQ{URL,Queue}` instead of
`Server{Port}`.
