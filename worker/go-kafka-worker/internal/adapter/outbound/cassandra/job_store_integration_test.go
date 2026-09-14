//go:build integration

package cassandra_test

import (
	"context"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/makifbaysal/boilerplates/worker/go-kafka-worker/internal/adapter/outbound/cassandra"
	"github.com/makifbaysal/boilerplates/worker/go-kafka-worker/internal/core/job"
	cassandraclient "github.com/makifbaysal/boilerplates/worker/go-kafka-worker/pkg/cassandra"
)

// Run via docker-compose.dev.yml — see that file for the exact command.
func TestJobStore_RecordGetList(t *testing.T) {
	hosts := os.Getenv("CASSANDRA_TEST_HOSTS")
	if hosts == "" {
		t.Skip("CASSANDRA_TEST_HOSTS not set, skipping cassandra integration test")
	}

	session, err := cassandraclient.Connect(cassandraclient.Config{
		Hosts:    strings.Split(hosts, ","),
		Keyspace: "worker_test",
	})
	require.NoError(t, err)
	defer session.Close()

	require.NoError(t, session.Query(`CREATE TABLE IF NOT EXISTS jobs (
		id text PRIMARY KEY, task_id text, type text, payload text, received_at timestamp, processed_at timestamp)`).Exec())

	store := cassandra.NewJobStore(session)
	ctx := context.Background()
	id := uuid.NewString()
	t.Cleanup(func() { _ = session.Query(`DELETE FROM jobs WHERE id = ?`, id).Exec() })

	now := time.Now().UTC().Truncate(time.Millisecond)
	in := &job.Job{ID: id, TaskID: "task-1", Type: "task.created", Payload: "{}", ReceivedAt: now, ProcessedAt: now}
	require.NoError(t, store.Record(ctx, in))

	got, err := store.Get(ctx, id)
	require.NoError(t, err)
	require.Equal(t, "task-1", got.TaskID)
}
