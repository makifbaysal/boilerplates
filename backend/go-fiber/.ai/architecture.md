# Architecture

Hexagonal (ports & adapters). The core knows nothing about HTTP, gRPC,
GraphQL, or any specific database — every framework touches it through an
interface defined *by* the core.

```
cmd/
  api/main.go              wiring only: config -> logger -> repo -> service -> rest adapter -> Listen
  grpc/main.go               same wiring, mounts the gRPC adapter instead
  graphql/main.go             same wiring, mounts the GraphQL adapter instead
internal/
  core/task/                  the hexagon: no fiber/grpc/gqlgen/db-driver imports, ever
    task.go                    domain struct (no json/koanf/db tags — adapters own their own DTOs)
    ports.go                    Repository (outbound port) + UseCase (inbound port) + sentinel errors
    service.go                  UseCase impl — the only thing that calls Repository
    mocks/                      mockery output for Repository, generated — never edit by hand
  adapter/
    inbound/                    driving adapters: protocol -> UseCase call
      rest/                      fiber v3 (see .ai/common-tasks.md)
      grpc/                      google.golang.org/grpc + protoc-gen-go
      graphql/                    gqlgen
    outbound/                   driven adapters: Repository port -> real storage
      memory/                    default for cmd/api, and what adapter tests run against
      postgres/, mongo/, cassandra/   build on pkg/{postgres,mongo,cassandra}
  platform/
    logger/                    log/slog (stdlib) + fiber middleware + grpc unary interceptor
    config/                     koanf loader, see "Config" below
pkg/
  postgres/, mongo/, cassandra/   reusable, app-agnostic client constructors (connect, ping,
                                   close) — no task-specific code. Import these from an
                                   outbound adapter; never import an outbound adapter from here.
  redis/, elasticsearch/          same idea, but not wired to a Repository — redis is for
                                   cache/rate-limit, elasticsearch for search/analytics indexes,
                                   neither is a sane system of record for task.Task.
configs/
  config.dev.yaml, config.staging.yaml, config.prod.yaml
```

Generated code (`internal/adapter/inbound/grpc/taskpb/`,
`internal/adapter/inbound/graphql/{generated,model}/`,
`internal/core/task/mocks/`) is committed to the repo — CI builds/tests
what's checked in, it does not run protoc/gqlgen/mockery. Regenerate
locally after changing a `.proto`/`.graphqls`/interface with `make
proto` / `make gqlgen` / `make mocks` and commit the diff.

## Layering rule

`adapter/inbound/*` never touches `Repository` or storage directly —
always through `task.UseCase`. `core/task` never imports a framework
package (no `fiber`, `grpc`, `gqlgen`, no DB driver) — that boundary is
what makes `service.go` unit-testable against a mocked `Repository`
without spinning up HTTP/gRPC/a real database, and what lets three
different transports (rest/grpc/graphql) share one implementation.

## Adding a new resource

Copy the `internal/core/task` shape (task.go → ports.go → service.go),
then add it to whichever inbound adapter(s) you need (a handler in
`rest`, a service method in `grpc`, a resolver in `graphql`) and an
outbound implementation (start with `memory`). Add the new `Repository`
to `.mockery.yaml` and run `make mocks`.

## Adding a new transport to an existing resource

The use case (`task.UseCase`) doesn't change — only add a new
`adapter/inbound/<transport>` that translates that protocol into calls
against the same `task.UseCase` interface `rest` already uses. See
`.ai/common-tasks.md`.

## Error handling

Domain errors are sentinel values (`errors.New`, checked with
`errors.Is`) defined in `core/task/ports.go`. Each inbound adapter maps
them to its own transport's error shape locally (`mapError` in `rest`,
a `status.Error` in `grpc`, a resolver error in `graphql`) — the core
never returns a transport-specific error type.

## Logging

One `*slog.Logger` (stdlib `log/slog`, no extra dependency) built in
`cmd/*/main.go` via `platform/logger.NewDynamic`. Inbound adapters attach
a request-scoped logger (with `request_id`) to the request's
`context.Context` — `logger.FiberMiddleware` for REST,
`logger.UnaryServerInterceptor` for gRPC. Anything below the adapter
layer that wants to log pulls it back out with `logger.FromContext(ctx)`
instead of taking `*slog.Logger` as a constructor argument, so the core
stays framework-free.

## Config

`platform/config` loads `configs/config.<APP_ENV>.yaml` via koanf
(`APP_ENV` defaults to `dev` for local `go run`/tests), then layers
`APP_`-prefixed environment variables on top as overrides
(`APP_SERVER_PORT` → `server.port`, etc — env always wins). `Loader.Watch`
follows the file with fsnotify and hot-swaps the in-memory config on
change; `cmd/api/main.go` wires that into a `slog.LevelVar` so editing
the log level in the config file takes effect without a restart.

**The Docker image does not default `APP_ENV`.** One image is built and
promoted across dev → staging → prod, so baking a default (e.g. `prod`)
would risk staging silently loading prod config the day someone forgets
to pass it. The deployer sets `APP_ENV` explicitly — `docker run -e
APP_ENV=staging`, the k8s manifest's `env:`, docker-compose's
`environment:`, etc. `CONFIG_DIR` is set in the image to `/configs`
where the demo files are copied to; override it if you mount configs
from elsewhere.
