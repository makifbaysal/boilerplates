# Coding standards

Same rules as `backend/go-fiber/.ai/coding-standards.md`; the
worker-specific notes are below.

## Layering (see architecture.md)

- `core/job` imports nothing outside the stdlib + `github.com/google/uuid`.
  No `franz-go`, no `amqp091-go`, no DB driver. Ever.
- Ports (`Store`, `Processor`) are defined *by the core*.
- An inbound consumer depends on `job.Processor` (the interface), never
  on `*job.Service` (the concrete type).

## Errors

- Domain errors are sentinel values next to the interface method that
  returns them (`job.ErrInvalidPayload`, `job.ErrNotFound` in `ports.go`),
  checked with `errors.Is`.
- Each inbound consumer's `HandleMessage` returns the raw error (decode
  error or `job.Processor`'s error) — `Run` is what decides the delivery
  policy (currently: log and Ack/commit regardless, see architecture.md).

## Dependency injection

- Constructor injection only. `cmd/{kafka,rabbitmq}/main.go` is the only
  place allowed to construct concrete types and wire them into
  interfaces — same rule as go-fiber.

## Style

- `gofmt`+`go vet` clean, `golangci-lint run` clean (`make lint`).
- No comments that restate the code.
- One exported type's worth of responsibility per file.
- Table-driven tests, small composable functions — see `.ai/testing.md`.

## Config & secrets

- Nothing reads `os.Getenv` outside `cmd/*/main.go` and
  `internal/platform/config`.
- No secrets in `configs/*.yaml` — real secrets come from `APP_`-prefixed
  env vars at deploy time.
