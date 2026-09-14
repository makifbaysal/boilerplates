// Package cassandra implements job.Store against gocql. The session
// this type takes as a constructor argument comes from pkg/cassandra.
// See schema.cql for the table this expects.
package cassandra

import (
	"context"
	"errors"

	"github.com/gocql/gocql"

	"github.com/makifbaysal/boilerplates/worker/go-kafka-worker/internal/core/job"
)

type JobStore struct {
	session *gocql.Session
}

var _ job.Store = (*JobStore)(nil)

func NewJobStore(session *gocql.Session) *JobStore {
	return &JobStore{session: session}
}

func (s *JobStore) Record(ctx context.Context, j *job.Job) error {
	return s.session.Query(
		`INSERT INTO jobs (id, task_id, type, payload, received_at, processed_at) VALUES (?, ?, ?, ?, ?, ?)`,
		j.ID, j.TaskID, j.Type, j.Payload, j.ReceivedAt, j.ProcessedAt,
	).WithContext(ctx).Exec()
}

func (s *JobStore) Get(ctx context.Context, id string) (*job.Job, error) {
	var j job.Job
	err := s.session.Query(
		`SELECT id, task_id, type, payload, received_at, processed_at FROM jobs WHERE id = ?`, id,
	).WithContext(ctx).Scan(&j.ID, &j.TaskID, &j.Type, &j.Payload, &j.ReceivedAt, &j.ProcessedAt)
	if errors.Is(err, gocql.ErrNotFound) {
		return nil, job.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &j, nil
}

func (s *JobStore) List(ctx context.Context) ([]*job.Job, error) {
	iter := s.session.Query(`SELECT id, task_id, type, payload, received_at, processed_at FROM jobs`).WithContext(ctx).Iter()

	out := make([]*job.Job, 0)
	var j job.Job
	for iter.Scan(&j.ID, &j.TaskID, &j.Type, &j.Payload, &j.ReceivedAt, &j.ProcessedAt) {
		cp := j
		out = append(out, &cp)
	}
	if err := iter.Close(); err != nil {
		return nil, err
	}
	return out, nil
}
