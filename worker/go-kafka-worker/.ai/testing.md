# Testing

Same conventions as `backend/go-fiber/.ai/testing.md` — mockery for
mocking ports, testify suite + table-driven subtests, real-implementation
adapter tests over mocked-core tests where the real thing is cheap. The
worker-specific pattern is the consumer split below.

## Consumer tests: HandleMessage vs. Run

Every consumer (`kafka.Consumer`, `rabbitmq.Consumer`) splits in two:

- `HandleMessage(ctx, source, value []byte) error` — decode + call
  `job.Processor`. Pure, dependency-free beyond `job.Processor`, and
  what `consumer_test.go` exercises directly with a real
  `memory.NewJobStore()` (table-driven: valid message, invalid JSON,
  invalid payload) — no broker needed.
- `Run(ctx, ...)` — the polling/channel loop around `HandleMessage`.
  Needs a live broker, so it's covered by
  `consumer_integration_test.go` instead (`//go:build integration`,
  skips unless `KAFKA_TEST_BROKERS` / `RABBITMQ_TEST_URL` is set — run
  via `docker-compose.dev.yml`, produces one message and asserts it
  lands in the store).

Keep this split when adding a third consumer type: a pure
`HandleMessage` you can unit test, a thin `Run` you integration-test.

## Integration tests against real infra

`internal/adapter/inbound/{kafka,rabbitmq}/consumer_integration_test.go`
and `internal/adapter/outbound/{postgres,mongo,cassandra}/*_integration_test.go`
are all gated behind `//go:build integration` and skip unless their env
var is set. CI does not run these — it builds/vets/tests everything
else, which is enough to catch a broken adapter signature; only real
infra catches a broken query or a broker-wiring bug.

## Coverage

No hard percentage gate, but `core/job/service.go` and each consumer's
`HandleMessage` should have positive-path *and* negative-path (invalid
payload, decode failure) coverage.
