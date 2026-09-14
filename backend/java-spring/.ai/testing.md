# Testing

## Mocking: Mockito, always

Never hand-write a mock. `@ExtendWith(MockitoExtension.class)` +
`@Mock`/`@InjectMocks` — mock only *ports* (`TaskRepositoryPort`), never
a concrete adapter, and never mock `TaskUseCase` just to test an inbound
adapter (see "adapter tests" below). This is the direct equivalent of
go-fiber's mockery — Mockito needs no code generation step, the mock is
created at test-run time.

## Structure: `@Nested` + `@ParameterizedTest`

One test class per unit under test. Group related cases with `@Nested`
classes (`TaskServiceTest.Create`, `.SetDone`) — the JUnit 5 equivalent
of go-fiber's `testify/suite`, sharing a `@BeforeEach`-built fixture.
Table-drive cases with `@ParameterizedTest` + `@MethodSource`/`@CsvSource`
rather than copy-pasted near-duplicate `@Test` methods — see
`TaskServiceTest.Create.create(...)` for the shape to copy.

## Adapter tests: real core, not mocks

`TaskControllerTest` (`@WebMvcTest`), `TaskGrpcServiceTest` (in-process
gRPC via `io.grpc.inprocess`), and `TaskGraphQLControllerTest`
(`@GraphQlTest`) all run the real adapter against a real
`InMemoryTaskRepository` — no mocks, importing `TaskService` +
`InMemoryTaskRepository` directly rather than stubbing `TaskUseCase`.
Same reasoning as go-fiber: the use case is cheap and side-effect-free
enough that an honest end-to-end adapter test beats a mocked one, and it
exercises the real (de)serialization each transport does.

## Integration tests against real databases: Testcontainers

`*IT.java` under `adapter.outbound.{postgres,mongo,cassandra}` use
`@Testcontainers` + `@SpringBootTest(properties = "grpc.server.port=-1")`
+ `@ActiveProfiles("<db>")`, with `@DynamicPropertySource` wiring the
container's real connection details in. Maven's **failsafe** plugin
picks up `*IT.java` only on `mvn verify` — plain `mvn test` (surefire)
does not run them, same skip-by-default shape as go-fiber's
`//go:build integration` tag. Needs a real Docker daemon; no manual env
var required (Testcontainers manages the container lifecycle itself,
unlike go-fiber's `docker-compose.dev.yml` + env-var-gated skip).

## Regression guard: `SmokeTest`

A bare `@SpringBootTest` with no DB profile active, asserting the
context loads. Exists specifically to catch the
all-three-DB-autoconfigurations-fight-each-other failure documented in
`.ai/architecture.md#choosing-a-database-profiles-not-constructor-swaps`
— don't delete it even though it "doesn't test anything."

## Coverage

No hard percentage gate, but `TaskService` and each inbound adapter
should have positive-path *and* negative-path (invalid input,
not-found) coverage — see any existing test class for the shape to copy.
