// Package postgres implements task.Repository against pgx/pgxpool. See
// schema.sql for the table this expects; pkg/postgres builds the pool
// this type takes as a constructor argument.
package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/makifbaysal/boilerplates/backend/go-fiber/internal/core/task"
)

type TaskRepository struct {
	pool *pgxpool.Pool
}

var _ task.Repository = (*TaskRepository)(nil)

func NewTaskRepository(pool *pgxpool.Pool) *TaskRepository {
	return &TaskRepository{pool: pool}
}

func (r *TaskRepository) Create(ctx context.Context, t *task.Task) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO tasks (id, title, done, created_at, updated_at) VALUES ($1, $2, $3, $4, $5)`,
		t.ID, t.Title, t.Done, t.CreatedAt, t.UpdatedAt)
	return err
}

func (r *TaskRepository) Get(ctx context.Context, id string) (*task.Task, error) {
	var t task.Task
	err := r.pool.QueryRow(ctx,
		`SELECT id, title, done, created_at, updated_at FROM tasks WHERE id = $1`, id,
	).Scan(&t.ID, &t.Title, &t.Done, &t.CreatedAt, &t.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, task.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *TaskRepository) List(ctx context.Context) ([]*task.Task, error) {
	rows, err := r.pool.Query(ctx, `SELECT id, title, done, created_at, updated_at FROM tasks ORDER BY created_at`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]*task.Task, 0)
	for rows.Next() {
		var t task.Task
		if err := rows.Scan(&t.ID, &t.Title, &t.Done, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, err
		}
		out = append(out, &t)
	}
	return out, rows.Err()
}

func (r *TaskRepository) Update(ctx context.Context, t *task.Task) error {
	tag, err := r.pool.Exec(ctx,
		`UPDATE tasks SET title = $2, done = $3, updated_at = $4 WHERE id = $1`,
		t.ID, t.Title, t.Done, t.UpdatedAt)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return task.ErrNotFound
	}
	return nil
}

func (r *TaskRepository) Delete(ctx context.Context, id string) error {
	tag, err := r.pool.Exec(ctx, `DELETE FROM tasks WHERE id = $1`, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return task.ErrNotFound
	}
	return nil
}
