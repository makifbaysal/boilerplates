.PHONY: run-api run-grpc run-graphql build test test-integration lint mocks proto gqlgen tidy

run-api:
	go run ./cmd/api

run-grpc:
	go run ./cmd/grpc

run-graphql:
	go run ./cmd/graphql

build:
	go build -o bin/api ./cmd/api
	go build -o bin/grpc ./cmd/grpc
	go build -o bin/graphql ./cmd/graphql

test:
	go test ./... -race -cover

# Requires docker-compose.dev.yml services running — see that file.
test-integration:
	go test ./internal/adapter/outbound/... -tags=integration -race -v

lint:
	go vet ./...
	@if command -v golangci-lint >/dev/null 2>&1; then golangci-lint run ./...; \
	else go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.12.2 run ./...; fi

mocks:
	go run github.com/vektra/mockery/v2@v2.53.6

# Requires protoc (system package, e.g. `brew install protobuf`) plus
# protoc-gen-go/protoc-gen-go-grpc (go install ...@latest, once).
proto:
	protoc -I internal/adapter/inbound/grpc/proto \
		--go_out=. --go_opt=module=github.com/makifbaysal/boilerplates/backend/go-fiber \
		--go-grpc_out=. --go-grpc_opt=module=github.com/makifbaysal/boilerplates/backend/go-fiber \
		internal/adapter/inbound/grpc/proto/task.proto

gqlgen:
	go run github.com/99designs/gqlgen generate

tidy:
	go mod tidy
