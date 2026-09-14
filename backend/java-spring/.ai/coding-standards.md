# Coding standards

Same principles as `backend/go-fiber/.ai/coding-standards.md`; Java/Spring
specifics below.

## Layering

- `core.task` has no Spring/grpc/gqlgen/JPA imports beyond `@Service` on
  `TaskService` itself (wiring convenience, not a real dependency —
  everything the class *does* is plain Java).
- Adapters depend on `TaskUseCase`/`TaskRepositoryPort` (the ports),
  never on `TaskService` or a concrete outbound adapter type directly.

## Errors

- Domain errors are unchecked exceptions (`TaskNotFoundException`,
  `InvalidTaskTitleException`) in `core.task`, checked by type.
- Each inbound adapter maps them to its own transport's error shape
  once: `RestExceptionHandler` (`@RestControllerAdvice`), `mapError` in
  `TaskGrpcService`, `GraphQLExceptionResolver`. Never catch a domain
  exception deep inside a controller/service method to translate it —
  that logic lives in exactly one of these three places per transport.

## Dependency injection

- Constructor injection only, no field `@Autowired`. `@Profile`-gated
  outbound adapters are how the app chooses a database — see
  architecture.md — nothing else should branch on environment.

## Style

- Lombok (`@Getter`, `@NoArgsConstructor`, `@AllArgsConstructor`) on
  adapter-local entities/documents (JPA/Mongo/Cassandra) to cut
  boilerplate; never on `core.task.Task` itself (it's a record).
- No comments that restate the code — same bar as go-fiber.
- One exported type's worth of responsibility per file.
- JUnit 5 `@Nested` + `@ParameterizedTest` over deeply nested
  conditionals in tests — see `.ai/testing.md`.

## Config & secrets

- Nothing reads `System.getenv`/`@Value` for infra config outside
  `application*.yml` and `LiveLogLevelWatcher` (which explicitly needs
  the raw env var for its file-watch path resolution).
- No secrets in `application-*.yml` — those are demo/dev defaults
  committed to git. Real secrets come from env vars
  (`SPRING_DATASOURCE_PASSWORD` etc, Spring's relaxed binding maps them
  automatically) at deploy time.
