# Common tasks

## Add a new endpoint / RPC / query

Same shape across all three transports: add the method to
`core.task.TaskUseCase`/`TaskService` first, then add the
transport-specific method that calls it —
`TaskResource` (`@GET`/`@POST`/... + `@Path`) for REST,
`TaskGrpcService` (edit `src/main/proto/task.proto` first, Quarkus
auto-compiles it — nothing to check in or regenerate manually) for
gRPC, `TaskGraphQLApi` (`@Query`/`@Mutation`, remember `@Name(...)` if
the method name would otherwise get bean-property-stripped, e.g. a
`setXxx` method) for GraphQL. Add a test alongside each.

## Add a brand-new resource

Copy `core.task`'s shape into `core.<resource>`, then an outbound
implementation starting with an in-memory one, then whichever inbound
adapter(s) you need. If a second implementation of the port exists,
both need `@Typed(OwnClass.class)` and a selector — see
architecture.md.

## Swap the outbound store

Set `app.datastore=postgres`/`mongo`/`cassandra` (env var
`APP_DATASTORE` also works via MicroProfile Config's relaxed binding).
`TaskRepositorySelector` handles the rest — no code change.
