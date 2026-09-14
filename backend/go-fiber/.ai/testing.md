# Testing

## Mocking: mockery, always

Never hand-write a mock. Every interface that needs one is registered in
`.mockery.yaml`; running `make mocks` (re)generates it under
`<package>/mocks/`. Mock only *ports* (`task.Repository`) — never mock a
concrete adapter type, and never mock `task.UseCase` just to test an
inbound adapter (see "adapter tests" below for why).

```go
repo := mocks.NewRepository(s.T())          // registers t.Cleanup(AssertExpectations)
repo.EXPECT().Get(mock.Anything, "abc").Return(existing, nil).Once()
```

## Structure: testify suite + table-driven subtests

One `testify/suite.Suite` per unit under test (one per `_test.go`,
named `<Thing>Suite`). Shared fixtures (a fresh mock, a fresh adapter)
go in `SetupTest`, which testify runs before every `Test*` method — see
`internal/core/task/service_test.go`.

Inside each `Test*` method, table-drive the cases:

```go
func (s *ServiceSuite) TestCreate() {
    cases := []struct {
        name      string
        title     string
        setupMock func()
        wantErr   error
    }{ /* ... */ }

    for _, tc := range cases {
        s.Run(tc.name, func() {
            s.SetupTest() // fresh mock per case — SetupTest only fires once per Test* otherwise
            tc.setupMock()
            // ...assert...
        })
    }
}
```

Call `s.SetupTest()` explicitly at the top of each `s.Run` when a table
test needs isolated mock expectations per case (testify's `SetupTest`
hook fires once per `Test*` method, not per subtest).

## Adapter tests: real core, not mocks

`rest.HandlerSuite`, `grpc.ServerSuite`, and `graphql.ResolverSuite` all
run the real generated/wired server against a real
`memory.NewTaskRepository()` — no mocks. The use case here is cheap and
side-effect-free enough that an honest end-to-end adapter test is both
simpler and more trustworthy than mocking `task.UseCase`; it also
exercises the actual (de)serialization each transport does. Reach for a
mocked `Repository` at the `core/task` service-test level instead, where
you want to assert exact calls/args without a real store.

## Integration tests against real databases

`internal/adapter/outbound/{postgres,mongo,cassandra}/*_integration_test.go`
are gated behind `//go:build integration` and skip unless their
`*_TEST_DSN`/`*_TEST_URI`/`*_TEST_HOSTS` env var is set. Run them via
`docker-compose.dev.yml` (see that file for the exact command). CI does
not run these — it builds/vets/tests everything else, which is enough
to catch a broken adapter signature; only a real database catches a
broken query.

## Coverage

No hard percentage gate, but a resource's `service.go` and its inbound
adapter handler/resolver should each have positive-path *and*
negative-path (invalid input, not-found) coverage — see any existing
`Test*` for the shape to copy.
