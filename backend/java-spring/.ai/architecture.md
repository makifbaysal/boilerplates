# Architecture

Hexagonal (ports & adapters), same shape as `backend/go-fiber` — the
core knows nothing about Spring MVC, gRPC, GraphQL, or any specific
database.

```
core/task/                    the hexagon: plain Java, no Spring/grpc/gqlgen/JPA imports
  Task.java                     domain record
  TaskRepositoryPort.java        outbound port
  TaskUseCase.java                 inbound port
  TaskService.java                  @Service implementing TaskUseCase — only thing calling the port
  TaskNotFoundException.java, InvalidTaskTitleException.java   sentinel domain exceptions
adapter/
  inbound/
    rest/           Spring MVC — TaskController + dto/ + RestExceptionHandler
    grpc/            net.devh grpc-server-spring-boot-starter — TaskGrpcService, proto/task.proto
    graphql/          spring-boot-starter-graphql — TaskGraphQLController + GraphQLExceptionResolver
  outbound/
    postgres/, mongo/, cassandra/    Spring Data {JPA,MongoDB,Cassandra}, one @Profile each
platform/
  logging/          RequestLoggingFilter (REST) + GrpcLoggingInterceptor (gRPC) + LiveLogLevelWatcher
```

## One process, not three binaries

Unlike go-fiber (separate `cmd/{api,grpc,graphql}` binaries), Spring Boot
runs everything in one JVM process: REST (`/tasks`) and GraphQL
(`/graphql`) share the same embedded Tomcat on port 8080; the gRPC
server (net.devh) runs in-process on its own port 9090 alongside it.
One Dockerfile, two exposed ports.

## Layering rule

Same as go-fiber: `adapter.inbound.*` never touches `TaskRepositoryPort`
directly — always through `TaskUseCase`. `core.task` imports nothing
Spring-specific (`TaskService` is `@Service`-annotated for wiring
convenience only — it has no other Spring dependency, and the interface
it implements, `TaskUseCase`, is plain Java).

## Choosing a database: profiles, not constructor swaps

Go swaps the outbound `Repository` implementation in `main.go`. Spring's
idiom is `@Profile` on each adapter (`@Profile("postgres")` etc, default
`InMemoryTaskRepository` active when none of those profiles are set) —
activate with `SPRING_PROFILES_ACTIVE=postgres` (combine with an env
profile: `SPRING_PROFILES_ACTIVE=dev,postgres`).

**Why every DB autoconfiguration is excluded by default in
`application.yml`:** postgres/mongo/cassandra starters are all on one
classpath (so any of the three profiles can activate), and Spring Boot
autoconfigures whatever it finds on the classpath regardless of *your*
custom `@Profile` annotations — with all three present it tries to wire
up all three DataSources/clients at once and refuses to start. Each
`application-{postgres,mongo,cassandra}.yml` re-excludes only the
*other* two, which is what actually turns its own DB on.
`SmokeTest` guards against this regressing.

## Logging

`RequestLoggingFilter` (REST) and `GrpcLoggingInterceptor` (gRPC) each
put a `request_id` in SLF4J's MDC for the duration of the
request/call — Java's request-scoped-logger mechanism, playing the same
role go-fiber's `logger.WithContext(ctx, ...)` does. Anything logged
downstream during that request picks it up automatically via Logback's
pattern layout, no explicit passing required.

## Config

`SPRING_PROFILES_ACTIVE` combines an env profile (`dev`/`staging`/`prod`
— sets `logging.level.root`) with a DB profile
(`postgres`/`mongo`/`cassandra` — sets connection properties + the
autoconfigure-exclude override above; omit for the in-memory default).

`LiveLogLevelWatcher` watches the active env profile's
`application-<profile>.yml` **on the real filesystem** (`CONFIG_DIR` if
set, else `src/main/resources` for local `mvn spring-boot:run`) and
applies a changed `logging.level.root` without a restart — the Java
equivalent of go-fiber's `config.Loader.Watch` wired into a
`slog.LevelVar`. The Docker image copies the yml files to `/config` and
sets `SPRING_CONFIG_ADDITIONAL_LOCATION=file:/config/` so that external
copy is what Spring actually reads, not just what the watcher reads —
edit it live in a running container and both the watcher and (on
restart) Spring's own config pick it up. Mirrors go-fiber's "the image
does not default `SPRING_PROFILES_ACTIVE`" reasoning: one image
promoted across environments, deployer sets it explicitly.
