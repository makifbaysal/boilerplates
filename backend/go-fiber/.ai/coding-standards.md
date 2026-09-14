# Coding standards

## Layering (see architecture.md)

- `core/task` imports nothing outside the stdlib + `github.com/google/uuid`.
  No `fiber`, no `grpc`, no `gqlgen`, no DB driver. Ever.
- Ports (`Repository`, `UseCase`) are defined *by the core*, not by the
  adapters that implement/consume them — that's what makes them ports.
- An inbound adapter depends on `task.UseCase` (the interface), never on
  `*task.Service` (the concrete type).

## Errors

- Domain errors are sentinel values (`var ErrX = errors.New(...)`) next
  to the interface method that returns them, checked with `errors.Is`.
- Never a bare `errors.New("...")` at a call site for something a caller
  might need to branch on — if it's checkable, it's a sentinel (or a
  wrapped one via `fmt.Errorf("...: %w", err)`).
- Each inbound adapter maps domain errors to its own transport's error
  shape once, in one place (`mapError` in rest/grpc, a resolver error in
  graphql) — never scatter status-code decisions across handlers.

## Dependency injection

- Constructor injection only (`NewX(dep Y) *X`). No globals, no
  package-level `init()` wiring, no service locators.
- `cmd/*/main.go` is the only place allowed to construct concrete types
  (a `*memory.TaskRepository`, a `*pgxpool.Pool`, …) and wire them into
  interfaces. Nothing under `internal/core` or `internal/adapter` should
  reach for a concrete outbound type by name.

## Style

- `gofmt`+`go vet` clean, `golangci-lint run` clean (`make lint`) —
  non-negotiable, CI enforces it.
- No comments that restate the code. A comment earns its place only by
  explaining a non-obvious *why* (a workaround, an invariant, a protocol
  quirk) — see any `mapError`/`Watch`/`SetContext` call in this repo for
  the bar.
- Keep one exported type's worth of responsibility per file
  (`task_repository.go`, not `repository_and_helpers_and_types.go`).
- Prefer table-driven tests and small, composable functions over deeply
  nested conditionals — see `.ai/testing.md`.

## Config & secrets

- Nothing reads `os.Getenv` outside `cmd/*/main.go` and
  `internal/platform/config`. Business logic takes configuration as
  constructor arguments, not by reaching into the environment.
- No secrets in `configs/*.yaml` — those are demo/dev defaults committed
  to git. Real secrets come from `APP_`-prefixed env vars at deploy time
  (see architecture.md#config).
