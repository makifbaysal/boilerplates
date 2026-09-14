# Architecture

Hexagonal (ports & adapters), same shape as `backend/go-fiber` /
`backend/java-spring`.

```
core/task/                    plain Java, no Quarkus/grpc/smallrye/panache imports
adapter/
  inbound/
    rest/          Quarkus REST (JAX-RS) — TaskResource + RestExceptionMapper
    grpc/            quarkus-grpc — TaskGrpcService (proto/task.proto, auto-compiled, not checked in)
    graphql/          SmallRye GraphQL — TaskGraphQLApi (code-first, no schema.graphqls)
  outbound/
    postgres/, mongo/    Panache repository pattern (Hibernate ORM / MongoDB)
    cassandra/            plain DataStax driver — no maintained Quarkus extension exists
platform/
  logging/         RequestLoggingFilter (JAX-RS) + GrpcLoggingInterceptor
  datastore/        TaskRepositorySelector — see "Choosing a database" below
```

One process, like java-spring: REST+GraphQL share the same HTTP port
(8080), gRPC runs in-process on its own port (9090).

## Choosing a database: three real bugs, fixed

`app.datastore` (`memory`\|`postgres`\|`mongo`\|`cassandra`) picks the
active `TaskRepositoryPort`. Getting this working surfaced three
Quarkus-specific problems that don't exist in go-fiber or java-spring,
each caught by running the actual test suite, not by inspection:

1. **`@LookupIfProperty` doesn't gate normal `@Inject` sites.** It looks
   like Spring's `@Profile` but only affects programmatic
   `Arc.container().select(...)` lookups. With all four
   `TaskRepositoryPort` implementations `@ApplicationScoped` and
   `app.datastore` a runtime (not build-time) property, ArC can't
   exclude beans at build time and throws
   `AmbiguousResolutionException`. Fixed by adding
   `@Typed(OwnClass.class)` to each implementation (restricts its CDI
   bean types to itself, removing it from `TaskRepositoryPort`
   resolution) and adding `platform/datastore/TaskRepositorySelector`,
   a single `@Produces` method that's the *only* thing exposing
   `TaskRepositoryPort` — it injects `Instance<X>` for each candidate
   and only `.get()`s the selected one, so unselected adapters (and
   anything they'd eagerly connect to, like Cassandra's `CqlSession`)
   never get instantiated.
2. **Hibernate ORM activates from `@Entity` presence, not from
   `app.datastore`.** Even with the selector above, `quarkus-hibernate-orm-panache`
   still tries to build/validate a persistence unit at startup purely
   because `TaskPanacheEntity` is annotated `@Entity` — independent of
   whether postgres is actually selected. `quarkus.hibernate-orm.active: false`
   in `application.yml` turns that off by default;
   `PostgresTestResource` flips it back on (alongside a real
   Testcontainers Postgres) for just `PostgresTaskRepositoryIT`.
3. **`@QuarkusTestResource` isn't class-scoped by default.** All
   `@QuarkusTest` classes share one continuous-testing JVM, and Quarkus
   discovers `@QuarkusTestResource` annotations across the whole
   compiled test tree up front — without
   `restrictToAnnotatedClass = true`, an unrelated test (e.g. the
   GraphQL adapter test) could trip over `CassandraTestResource` trying
   to start a container it has nothing to do with. All three DB test
   resources set this explicitly.

Also disabled by default in `%test`: Quarkus Dev Services, which
otherwise auto-provisions a Testcontainers instance for *every*
DB/cache extension on the classpath (five containers, ~100s, for a test
that only needs the in-memory default).

## Errors

Same "map once, per transport" rule as the other backends:
`RestExceptionMapper` (JAX-RS `ExceptionMapper`), `mapError` inside
`TaskGrpcService`, and — GraphQL's own gotcha — SmallRye masks any
unmapped `RuntimeException` as a generic `"System error"` by default
(a safe default). `mp.graphql.showErrorMessage` in `application.yml`
allow-lists `TaskNotFoundException`/`InvalidTaskTitleException` so
their real message reaches the client; anything else stays masked.
(Also caught live: `@Mutation` strips a leading `set` from a method
name as if it were a bean setter — `setTaskDone` silently became field
`taskDone` until `@Name("setTaskDone")` forced it.)

## Dynamic config

`mvn quarkus:dev` live-reloads config *and code* on save — Quarkus's
own answer to the "edit config, see it take effect without a restart"
requirement go-fiber (koanf+fsnotify) and java-spring (WatchService)
implement by hand. That covers development; for a genuinely dynamic
*production* value the idiomatic Quarkus mechanism is a custom
MicroProfile `ConfigSource` that polls its backing store — not
implemented here to keep this boilerplate's scope honest, but that's
the extension point if a real project needs it.
