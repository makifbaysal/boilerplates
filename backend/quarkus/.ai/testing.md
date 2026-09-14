# Testing

Same conventions as the other backends: real adapter over real
in-memory core (not mocks) for REST/gRPC/GraphQL tests; Mockito +
JUnit5 `@Nested`/`@ParameterizedTest` for `core.task`'s own tests.

- `@QuarkusTest` boots the real app per test class (fast — Quarkus is
  built for this). Dev Services are disabled in `%test`
  (see architecture.md) so this stays fast even with 5 DB/cache
  extensions on the classpath.
- `SmokeTest` guards against the CDI-ambiguity/eager-Hibernate-ORM
  regression documented in architecture.md — don't delete it.
- DB integration tests (`*IT.java`) use `@QuarkusTestResource(..., restrictToAnnotatedClass = true)`
  backed by Testcontainers — **the `restrictToAnnotatedClass` flag is
  required**, not optional, or an unrelated test can trip over a
  container it has nothing to do with (see architecture.md). Run via
  `mvn verify` (needs Docker); `skipITs` is `false` here (Quarkus's
  generated default is `true`, meant for `-Dnative` only).
