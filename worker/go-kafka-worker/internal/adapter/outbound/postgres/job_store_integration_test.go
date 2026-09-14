//go:build integration

package postgres_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/makifbaysal/boilerplates/worker/go-kafka-worker/internal/adapter/outbound/postgres"
	"github.com/makifbaysal/boilerplates/worker/go-kafka-worker/internal/core/job"
	pgclient "github.com/makifbaysal/boilerplates/worker/go-kafka-worker/pkg/postgres"
)

// Run via docker-compose.dev.yml — see that file for the exact command.
func TestJobStore_RecordGetList(t *testing.T) {
	dsn := os.Getenv("POSTGRES_TEST_DSN")
	if dsn == "" {
		t.Skip("POSTGRES_TEST_DSN not set, skipping postgres integration test")
	}

	ctx := context.Background()
	pool, err := pgclient.Connect(ctx, pgclient.Config{DSN: dsn})
	require.NoError(t, err)
	defer pool.Close()

	_, err = pool.Exec(ctx, `CREATE TABLE IF NOT EXISTS jobs (
		id TEXT PRIMARY KEY, task_id TEXT NOT NULL, type TEXT NOT NULL, payload TEXT NOT NULL,
		received_at TIMESTAMPTZ NOT NULL, processed_at TIMESTAMPTZ NOT NULL)`)
	require.NoError(t, err)

	store := postgres.NewJobStore(pool)
	id := uuid.NewString()
	t.Cleanup(func() { _, _ = pool.Exec(ctx, `DELETE FROM jobs WHERE id = $1`, id) })

	now := time.Now().UTC().Truncate(time.Microsecond)
	in := &job.Job{ID: id, TaskID: "task-1", Type: "task.created", Payload: "{}", ReceivedAt: now, ProcessedAt: now}
	require.NoError(t, store.Record(ctx, in))

	got, err := store.Get(ctx, id)
	require.NoError(t, err)
	require.Equal(t, "task-1", got.TaskID)

	list, err := store.List(ctx)
	require.NoError(t, err)
	require.NotEmpty(t, list)
}
