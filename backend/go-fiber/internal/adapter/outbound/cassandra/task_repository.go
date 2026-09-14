// Package cassandra implements task.Repository against gocql. The
// session this type takes as a constructor argument comes from
// pkg/cassandra. See schema.cql for the table this expects.
package cassandra

import (
	"context"
	"errors"

	"github.com/gocql/gocql"

	"github.com/makifbaysal/boilerplates/backend/go-fiber/internal/core/task"
)

type TaskRepository struct {
	session *gocql.Session
}

var _ task.Repository = (*TaskRepository)(nil)

func NewTaskRepository(session *gocql.Session) *TaskRepository {
	return &TaskRepository{session: session}
}

func (r *TaskRepository) Create(ctx context.Context, t *task.Task) error {
	return r.session.Query(
		`INSERT INTO tasks (id, title, done, created_at, updated_at) VALUES (?, ?, ?, ?, ?)`,
		t.ID, t.Title, t.Done, t.CreatedAt, t.UpdatedAt,
	).WithContext(ctx).Exec()
}

func (r *TaskRepository) Get(ctx context.Context, id string) (*task.Task, error) {
	var t task.Task
	err := r.session.Query(
		`SELECT id, title, done, created_at, updated_at FROM tasks WHERE id = ?`, id,
	).WithContext(ctx).Scan(&t.ID, &t.Title, &t.Done, &t.CreatedAt, &t.UpdatedAt)
	if errors.Is(err, gocql.ErrNotFound) {
		return nil, task.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *TaskRepository) List(ctx context.Context) ([]*task.Task, error) {
	iter := r.session.Query(`SELECT id, title, done, created_at, updated_at FROM tasks`).WithContext(ctx).Iter()

	out := make([]*task.Task, 0)
	var t task.Task
	for iter.Scan(&t.ID, &t.Title, &t.Done, &t.CreatedAt, &t.UpdatedAt) {
		cp := t
		out = append(out, &cp)
	}
	if err := iter.Close(); err != nil {
		return nil, err
	}
	return out, nil
}

// Update is an upsert at the CQL level — Get first so the Repository
// port's "must already exist" contract holds the same way it does for
// postgres/mongo.
func (r *TaskRepository) Update(ctx context.Context, t *task.Task) error {
	if _, err := r.Get(ctx, t.ID); err != nil {
		return err
	}
	return r.session.Query(
		`UPDATE tasks SET title = ?, done = ?, updated_at = ? WHERE id = ?`,
		t.Title, t.Done, t.UpdatedAt, t.ID,
	).WithContext(ctx).Exec()
}

func (r *TaskRepository) Delete(ctx context.Context, id string) error {
	if _, err := r.Get(ctx, id); err != nil {
		return err
	}
	return r.session.Query(`DELETE FROM tasks WHERE id = ?`, id).WithContext(ctx).Exec()
}
