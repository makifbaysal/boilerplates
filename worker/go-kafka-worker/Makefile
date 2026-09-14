.PHONY: run-kafka run-rabbitmq build test test-integration lint mocks tidy

run-kafka:
	go run ./cmd/kafka

run-rabbitmq:
	go run ./cmd/rabbitmq

build:
	go build -o bin/kafka ./cmd/kafka
	go build -o bin/rabbitmq ./cmd/rabbitmq

test:
	go test ./... -race -cover

# Requires docker-compose.dev.yml services running — see that file.
test-integration:
	go test ./... -tags=integration -race -v

lint:
	go vet ./...
	@if command -v golangci-lint >/dev/null 2>&1; then golangci-lint run ./...; \
	else go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.12.2 run ./...; fi

mocks:
	go run github.com/vektra/mockery/v2@v2.53.6

tidy:
	go mod tidy
