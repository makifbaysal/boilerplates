//go:build integration

package mongo_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/makifbaysal/boilerplates/backend/go-fiber/internal/adapter/outbound/mongo"
	"github.com/makifbaysal/boilerplates/backend/go-fiber/internal/core/task"
	mongoclient "github.com/makifbaysal/boilerplates/backend/go-fiber/pkg/mongo"
)

// Run via docker-compose.dev.yml — see that file for the exact command.
func TestTaskRepository_CRUD(t *testing.T) {
	uri := os.Getenv("MONGO_TEST_URI")
	if uri == "" {
		t.Skip("MONGO_TEST_URI not set, skipping mongo integration test")
	}

	ctx := context.Background()
	client, err := mongoclient.Connect(ctx, mongoclient.Config{URI: uri})
	require.NoError(t, err)
	defer func() { _ = client.Disconnect(ctx) }()

	repo := mongo.NewTaskRepository(client, "boilerplate_test")
	id := uuid.NewString()
	t.Cleanup(func() {
		_, _ = client.Database("boilerplate_test").Collection("tasks").DeleteOne(ctx, map[string]any{"_id": id})
	})

	now := time.Now().UTC().Truncate(time.Millisecond)
	in := &task.Task{ID: id, Title: "buy milk", CreatedAt: now, UpdatedAt: now}
	require.NoError(t, repo.Create(ctx, in))

	got, err := repo.Get(ctx, id)
	require.NoError(t, err)
	require.Equal(t, "buy milk", got.Title)

	got.Done = true
	require.NoError(t, repo.Update(ctx, got))

	got, err = repo.Get(ctx, id)
	require.NoError(t, err)
	require.True(t, got.Done)

	list, err := repo.List(ctx)
	require.NoError(t, err)
	require.NotEmpty(t, list)

	require.NoError(t, repo.Delete(ctx, id))
	_, err = repo.Get(ctx, id)
	require.ErrorIs(t, err, task.ErrNotFound)
}
