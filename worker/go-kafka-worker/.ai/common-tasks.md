# Common tasks

## Add a new field to the message contract

Example: add `retry_count` to the `{"task_id","type"}` payload.

1. Add the field to `message` struct in **both**
   `internal/adapter/inbound/kafka/consumer.go` and
   `internal/adapter/inbound/rabbitmq/consumer.go` — they consume the
   same logical topic shape from two different brokers, keep them in
   sync.
2. If the core needs to act on it, add a parameter to
   `job.Processor.Process` (in `ports.go`) and thread it through
   `job.Service.Process` — same as any other core change.
3. Update both consumers' `HandleMessage` to pass the new value.
4. Add a case to both `consumer_test.go` table tests.

## Add a new message type / topic

Copy the `internal/core/job` shape if it's a genuinely different domain
(not just a new `Type` value on the same `Job`) — `job.go` → `ports.go`
→ `service.go` — then wire a new inbound consumer (or extend the
existing ones to also consume the new topic/queue) and an outbound
`Store` implementation starting with `memory`.

Often it's *not* a new domain — most new message types are just a new
`Type` string flowing through the existing `Job`/`Processor`/`Store`.
Only reach for a new core package when the new message needs different
fields entirely.

## Swap the outbound store

`cmd/{kafka,rabbitmq}/main.go` construct `memory.NewJobStore()` and hand
it to `job.NewService`. To use a real database instead: build the client
with `pkg/{postgres,mongo,cassandra}`, construct the matching
`internal/adapter/outbound/{postgres,mongo,cassandra}.NewJobStore(...)`,
and pass that to `job.NewService` instead — nothing else changes.

## Add a config field

Same as go-fiber: add it to the right struct in
`internal/platform/config/config.go` (`koanf:"..."` tag), add it to each
`configs/config.<env>.yaml`, read it via `cfgLoader.Current().<Field>`
in `cmd/*/main.go`. Automatically overridable via `APP_<PATH>` env vars.

## Change delivery/retry policy

Both consumers currently Ack/commit even on a `HandleMessage` error (see
`.ai/architecture.md#delivery-semantics`). To add retry or a
dead-letter queue: change the branch in `Run` (kafka:
`internal/adapter/inbound/kafka/consumer.go`; rabbitmq: `.../rabbitmq/consumer.go`)
where the error from `c.handle(...)` is currently discarded — that's the
one place delivery policy lives, `HandleMessage` itself stays unaware of
it.
