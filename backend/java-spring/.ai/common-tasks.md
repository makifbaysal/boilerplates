# Common tasks

## Add a new endpoint to an existing resource (REST)

Example: `PATCH /tasks/{id}/title` to rename a task.

1. **Core** (`core/task/TaskService.java`) — add the business method,
   using `get`+`update` on `TaskRepositoryPort` the same way `setDone`
   does. Add the method to `TaskUseCase` too.
2. **REST controller** (`adapter/inbound/rest/TaskController.java`) —
   add the `@PatchMapping("/{id}/title")` handler.
3. **Tests** — a `@Nested` case in `TaskServiceTest` (mocked
   `TaskRepositoryPort`) and a `@ParameterizedTest`/`@Test` in
   `TaskControllerTest` (real `InMemoryTaskRepository`).
4. `mvn test`.

## Add a new RPC (gRPC)

1. Edit `src/main/proto/task.proto` — add the message(s) and RPC.
2. `mvn compile` (the `protobuf-maven-plugin` regenerates
   `target/generated-sources/protobuf/**` automatically — nothing to
   commit, unlike go-fiber's checked-in `taskpb/`, since Maven
   regenerates on every build).
3. Implement the new method on `TaskGrpcService`, calling the same
   `TaskUseCase` method the REST controller and GraphQL controller call.
4. Add a case to `TaskGrpcServiceTest` (in-process, no real port).

## Add a new field/query/mutation (GraphQL)

1. Edit `src/main/resources/graphql/schema.graphqls`.
2. Add the `@QueryMapping`/`@MutationMapping` method to
   `TaskGraphQLController`, calling the same `TaskUseCase` method.
3. Add a case to `TaskGraphQLControllerTest` using `GraphQlTester`.

## Add a brand-new resource

Copy the `core.task` shape (`Task.java` → `*Port.java`/`*UseCase.java` →
`*Service.java`) into `core.<resource>`, then wire it into whichever
inbound adapter(s) you need and an outbound implementation, starting
with an in-memory one. No mockery-equivalent config file to update —
Mockito mocks any interface on demand, nothing to regenerate.

## Add a config field

1. Add it to the relevant `application*.yml` (base or a specific
   profile — see architecture.md for which file controls what).
2. Read it via constructor-injected `@Value`
   or a `@ConfigurationProperties` class — don't reach for
   `Environment`/`System.getenv` inside business code.
3. If it needs to react live like `logging.level.root` does, extend
   `LiveLogLevelWatcher` (or add a sibling watcher) rather than bolting
   more fields onto that one — keep each watcher single-purpose.

## Swap the outbound store

Set `SPRING_PROFILES_ACTIVE=postgres` (or `mongo`/`cassandra`) —
`@Profile` on each `TaskRepositoryPort` implementation handles the rest.
No code change. See architecture.md for why the *other* two DB
autoconfigurations must stay excluded when you do this.

## Add a new external dependency (cache, search index, etc.)

Add the Spring Data starter (e.g. `spring-boot-starter-data-redis`,
already present but unused — see go-fiber's `pkg/redis` note: it's for
caching/rate-limiting, not a `TaskRepositoryPort` implementation) and
inject the client (`RedisTemplate`, `ElasticsearchClient`) directly
where needed. It doesn't need to go through `TaskRepositoryPort` since
it isn't one.
