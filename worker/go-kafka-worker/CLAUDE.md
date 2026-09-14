# go-kafka-worker boilerplate

Go async-worker skeleton in hexagonal architecture: one `job` processing
use case consumed from two message brokers — Kafka (franz-go) and
RabbitMQ (amqp091-go) — with three outbound storage implementations
(Postgres, MongoDB, Cassandra, plus the in-memory default) borrowed
verbatim from `backend/go-fiber/pkg`, structured logging (`log/slog`)
with a message-handler middleware, and koanf-based config with live
reload.

Read [.ai/architecture.md](.ai/architecture.md) before changing
structure. Read [.ai/common-tasks.md](.ai/common-tasks.md) before
touching the message contract or adding a resource. See
[.ai/coding-standards.md](.ai/coding-standards.md) and
[.ai/testing.md](.ai/testing.md) for style/test conventions — same as
`backend/go-fiber`, adapted for consumers instead of HTTP/gRPC/GraphQL
handlers.

## Commands

```
make run-kafka / run-rabbitmq   # go run ./cmd/{kafka,rabbitmq}
make build                       # -> bin/{kafka,rabbitmq}
make test                        # go test ./... -race -cover
make test-integration            # real broker + DB adapters, needs docker-compose.dev.yml up
make lint                        # go vet + golangci-lint
make mocks                       # regenerate internal/core/job/mocks via mockery
```
