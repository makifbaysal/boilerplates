//go:build integration

package mongo_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/makifbaysal/boilerplates/worker/go-kafka-worker/internal/adapter/outbound/mongo"
	"github.com/makifbaysal/boilerplates/worker/go-kafka-worker/internal/core/job"
	mongoclient "github.com/makifbaysal/boilerplates/worker/go-kafka-worker/pkg/mongo"
)

// Run via docker-compose.dev.yml — see that file for the exact command.
func TestJobStore_RecordGetList(t *testing.T) {
	uri := os.Getenv("MONGO_TEST_URI")
	if uri == "" {
		t.Skip("MONGO_TEST_URI not set, skipping mongo integration test")
	}

	ctx := context.Background()
	client, err := mongoclient.Connect(ctx, mongoclient.Config{URI: uri})
	require.NoError(t, err)
	defer func() { _ = client.Disconnect(ctx) }()

	store := mongo.NewJobStore(client, "worker_test")
	id := uuid.NewString()
	t.Cleanup(func() {
		_, _ = client.Database("worker_test").Collection("jobs").DeleteOne(ctx, map[string]any{"_id": id})
	})

	now := time.Now().UTC().Truncate(time.Millisecond)
	in := &job.Job{ID: id, TaskID: "task-1", Type: "task.created", Payload: "{}", ReceivedAt: now, ProcessedAt: now}
	require.NoError(t, store.Record(ctx, in))

	got, err := store.Get(ctx, id)
	require.NoError(t, err)
	require.Equal(t, "task-1", got.TaskID)

	list, err := store.List(ctx)
	require.NoError(t, err)
	require.NotEmpty(t, list)
}
