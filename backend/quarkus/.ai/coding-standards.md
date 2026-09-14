# Coding standards

Same principles as `backend/go-fiber`/`backend/java-spring`.

- `core.task` has no Quarkus/CDI/grpc/graphql/panache imports.
- Adapters depend on `TaskUseCase`/`TaskRepositoryPort`, never a
  concrete type — `TaskRepositorySelector` is the one place that
  resolves which outbound adapter is live, see architecture.md.
- Domain errors are unchecked exceptions in `core.task`; each inbound
  adapter maps them once (`RestExceptionMapper`, `mapError` in
  `TaskGrpcService`, `mp.graphql.showErrorMessage` config).
- Constructor injection (`@Inject` on the constructor) only.
- `@Typed(OwnClass.class)` on any new outbound adapter that implements
  a port also implemented elsewhere — otherwise CDI's ambiguity kicks
  in the moment a second implementation exists (see architecture.md).
- Nothing outside `application*.yml` and `cmd`-equivalent wiring reads
  config directly — use `@ConfigProperty`.
