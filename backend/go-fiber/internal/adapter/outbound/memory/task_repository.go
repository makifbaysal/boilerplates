// Package memory implements the task.Repository port in-process. Used as
// the default in cmd/api/main.go and for adapter-level tests; swap it for
// adapter/outbound/postgres|mongo|cassandra behind the same interface.
package memory

import (
	"context"
	"sort"
	"sync"

	"github.com/makifbaysal/boilerplates/backend/go-fiber/internal/core/task"
)

type TaskRepository struct {
	mu    sync.RWMutex
	tasks map[string]*task.Task
}

func NewTaskRepository() *TaskRepository {
	return &TaskRepository{tasks: make(map[string]*task.Task)}
}

func (r *TaskRepository) Create(_ context.Context, t *task.Task) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.tasks[t.ID] = t
	return nil
}

func (r *TaskRepository) Get(_ context.Context, id string) (*task.Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	t, ok := r.tasks[id]
	if !ok {
		return nil, task.ErrNotFound
	}
	return t, nil
}

func (r *TaskRepository) List(_ context.Context) ([]*task.Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]*task.Task, 0, len(r.tasks))
	for _, t := range r.tasks {
		out = append(out, t)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].CreatedAt.Before(out[j].CreatedAt) })
	return out, nil
}

func (r *TaskRepository) Update(_ context.Context, t *task.Task) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.tasks[t.ID]; !ok {
		return task.ErrNotFound
	}
	r.tasks[t.ID] = t
	return nil
}

func (r *TaskRepository) Delete(_ context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.tasks[id]; !ok {
		return task.ErrNotFound
	}
	delete(r.tasks, id)
	return nil
}
