# go-fiber boilerplate

Go backend skeleton in hexagonal (ports & adapters) architecture: one
`task` CRUD use case exposed over three inbound transports — REST
(Fiber v3), gRPC, and GraphQL (gqlgen) — with three outbound storage
implementations (Postgres, MongoDB, Cassandra, plus the in-memory
default), structured logging (`log/slog`), and koanf-based config with
live reload.

Read [.ai/architecture.md](.ai/architecture.md) before changing structure —
it explains the layering and why it exists. Read
[.ai/common-tasks.md](.ai/common-tasks.md) before adding a new
endpoint/RPC/query or resource — copy that pattern instead of inventing
a new one. [.ai/coding-standards.md](.ai/coding-standards.md) and
[.ai/testing.md](.ai/testing.md) cover style and test conventions
(mockery + testify suite + table-driven).

## Commands

```
make run-api / run-grpc / run-graphql   # go run ./cmd/{api,grpc,graphql}
make build                               # -> bin/{api,grpc,graphql}
make test                                # go test ./... -race -cover
make test-integration                    # real DB adapters, needs docker-compose.dev.yml up
make lint                                # go vet + golangci-lint
make mocks                               # regenerate internal/core/task/mocks via mockery
make proto                               # regenerate internal/adapter/inbound/grpc/taskpb via protoc
make gqlgen                              # regenerate internal/adapter/inbound/graphql/{generated,model}
```

Ports (local dev): REST `:8080`, gRPC `:9090`, GraphQL `:8081`
(playground at `/`, endpoint at `/query`).
