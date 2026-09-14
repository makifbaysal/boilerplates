// Package postgres implements job.Store against pgx/pgxpool. See
// schema.sql for the table this expects; pkg/postgres builds the pool
// this type takes as a constructor argument.
package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/makifbaysal/boilerplates/worker/go-kafka-worker/internal/core/job"
)

type JobStore struct {
	pool *pgxpool.Pool
}

var _ job.Store = (*JobStore)(nil)

func NewJobStore(pool *pgxpool.Pool) *JobStore {
	return &JobStore{pool: pool}
}

func (s *JobStore) Record(ctx context.Context, j *job.Job) error {
	_, err := s.pool.Exec(ctx,
		`INSERT INTO jobs (id, task_id, type, payload, received_at, processed_at) VALUES ($1, $2, $3, $4, $5, $6)`,
		j.ID, j.TaskID, j.Type, j.Payload, j.ReceivedAt, j.ProcessedAt)
	return err
}

func (s *JobStore) Get(ctx context.Context, id string) (*job.Job, error) {
	var j job.Job
	err := s.pool.QueryRow(ctx,
		`SELECT id, task_id, type, payload, received_at, processed_at FROM jobs WHERE id = $1`, id,
	).Scan(&j.ID, &j.TaskID, &j.Type, &j.Payload, &j.ReceivedAt, &j.ProcessedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, job.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &j, nil
}

func (s *JobStore) List(ctx context.Context) ([]*job.Job, error) {
	rows, err := s.pool.Query(ctx, `SELECT id, task_id, type, payload, received_at, processed_at FROM jobs ORDER BY received_at`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]*job.Job, 0)
	for rows.Next() {
		var j job.Job
		if err := rows.Scan(&j.ID, &j.TaskID, &j.Type, &j.Payload, &j.ReceivedAt, &j.ProcessedAt); err != nil {
			return nil, err
		}
		out = append(out, &j)
	}
	return out, rows.Err()
}
