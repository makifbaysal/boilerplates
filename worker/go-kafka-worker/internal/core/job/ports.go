package job

import (
	"context"
	"errors"
)

var (
	ErrInvalidPayload = errors.New("job: task_id and type are required")
	ErrNotFound       = errors.New("job: not found")
)

// Store is the driven (outbound) port: where processed jobs get
// recorded. Implemented by internal/adapter/outbound/{memory,postgres,
// mongo,cassandra}. Mocked via mockery for service tests.
type Store interface {
	Record(ctx context.Context, j *Job) error
	Get(ctx context.Context, id string) (*Job, error)
	List(ctx context.Context) ([]*Job, error)
}

// Processor is the driving (inbound) port every consumer (kafka,
// rabbitmq) calls after decoding a message off the wire.
type Processor interface {
	Process(ctx context.Context, taskID, jobType, payload string) (*Job, error)
}
