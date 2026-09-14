package task

import (
	"context"
	"errors"
)

var (
	ErrNotFound     = errors.New("task: not found")
	ErrInvalidTitle = errors.New("task: title required")
)

// Repository is the driven (outbound) port: what the core needs from
// persistence. Implemented by internal/adapter/outbound/{memory,postgres,
// mongo,cassandra}. Mocked via mockery — see .mockery.yaml — for
// UseCase tests.
type Repository interface {
	Create(ctx context.Context, t *Task) error
	Get(ctx context.Context, id string) (*Task, error)
	List(ctx context.Context) ([]*Task, error)
	Update(ctx context.Context, t *Task) error
	Delete(ctx context.Context, id string) error
}

// UseCase is the driving (inbound) port: what every inbound adapter
// (rest, grpc, graphql) is allowed to call. Adapters depend on this
// interface, never on *Service directly, so the core can be swapped or
// wrapped (e.g. with a decorator) without touching adapter code.
type UseCase interface {
	Create(ctx context.Context, title string) (*Task, error)
	Get(ctx context.Context, id string) (*Task, error)
	List(ctx context.Context) ([]*Task, error)
	SetDone(ctx context.Context, id string, done bool) (*Task, error)
	Delete(ctx context.Context, id string) error
}
