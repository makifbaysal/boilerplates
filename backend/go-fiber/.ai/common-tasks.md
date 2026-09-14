# Common tasks

Concrete walkthroughs for the things people ask for repeatedly. Follow
the existing pattern in each file rather than restructuring. See
`.ai/testing.md` for how to test each of these, and `.ai/coding-standards.md`
for style rules.

## Add a new endpoint to an existing resource (REST)

Example: `PATCH /tasks/:id/title` to rename a task.

1. **Core** (`internal/core/task/service.go`) — add the business method,
   using `Get`+`Update` on `Repository` the same way `SetDone` does. No
   change needed in `ports.go` unless the resource needs a new
   *data-access* operation the `Repository` interface doesn't have yet.
2. **REST handler** (`internal/adapter/inbound/rest/handler.go`) — add
   the fiber method and register the route in `Register`:
   `router.Patch("/tasks/:id/title", h.rename)`.
3. **Tests** — a service-level case in `core/task/service_test.go`
   (mocked `Repository`, table-driven, see `.ai/testing.md`) and an
   adapter-level case in `rest/handler_test.go` (real in-memory repo).
4. `make test && make lint`.

## Add a new RPC (gRPC)

1. Edit `internal/adapter/inbound/grpc/proto/task.proto` — add the
   request/response messages and the RPC on `TaskService`.
2. `make proto` (runs protoc + protoc-gen-go + protoc-gen-go-grpc,
   regenerates `taskpb/*.go`) — commit the regenerated files.
3. Implement the new method on `*grpc.Server` in
   `internal/adapter/inbound/grpc/server.go`, calling the same
   `task.UseCase` method the REST handler and GraphQL resolver call.
4. Add a case to `grpc/server_test.go`'s `ServerSuite` (bufconn, no real
   port — see `.ai/testing.md`).

## Add a new field/query/mutation (GraphQL)

1. Edit `internal/adapter/inbound/graphql/schema.graphqls`.
2. `make gqlgen` (regenerates `generated/` and `model/`, and stubs any
   new resolver method into `schema.resolvers.go` as
   `panic("not implemented")`) — commit the regenerated files.
3. Replace the stub in `schema.resolvers.go`, calling `r.uc.<Method>`
   (the same `task.UseCase` every other transport calls) and mapping
   with `toModel`.
4. Add a case to `graphql/resolver_test.go`'s `ResolverSuite`.

## Add a brand-new resource

Copy the `internal/core/task` shape (`task.go` → `ports.go` →
`service.go`) into `internal/core/<resource>`, then wire it into
whichever inbound adapter(s) you need — a REST handler, a gRPC service
(new `.proto`, or a new RPC on an existing one), a GraphQL type/resolver
— and an outbound implementation, starting with a `memory.<Resource>Repository`.
Add the new `Repository` interface to `.mockery.yaml` and run `make mocks`.

## Add a config field

1. Add the field to the right struct in
   `internal/platform/config/config.go` (with a `koanf:"..."` tag).
2. Add it to each `configs/config.<env>.yaml` (or at least `dev`).
3. Read it via `cfgLoader.Current().<Field>` in `cmd/*/main.go` — or, if
   it needs to react to `Loader.Watch` live (like `Log.Level` already
   does for the `slog.LevelVar`), wire that in the `Watch` callback too.
4. It's automatically overridable via `APP_<PATH>` env vars — no extra
   code needed for that part.

## Swap the outbound store

`cmd/api/main.go` (and `cmd/grpc`, `cmd/graphql`) construct
`memory.NewTaskRepository()` and hand it to `task.NewService`. To use a
real database instead: build the client with `pkg/{postgres,mongo,cassandra}`,
construct the matching `internal/adapter/outbound/{postgres,mongo,cassandra}.NewTaskRepository(...)`,
and pass that to `task.NewService` instead — nothing else in any layer
changes, that's the point of the `Repository` port.

## Add a new external dependency (cache, search index, etc.)

If it's not a system-of-record for an existing resource (e.g. redis for
caching, elasticsearch for search — see `pkg/redis`, `pkg/elasticsearch`),
build the client in `pkg/<name>` (connect + health check, no task-specific
code), then use it directly from wherever needs it in `cmd/*/main.go` or
a new `internal/platform/<name>` helper — it doesn't need to go through
the `Repository` port since it isn't one.
