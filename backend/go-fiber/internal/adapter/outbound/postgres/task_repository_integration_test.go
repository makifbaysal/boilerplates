//go:build integration

package postgres_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/makifbaysal/boilerplates/backend/go-fiber/internal/adapter/outbound/postgres"
	"github.com/makifbaysal/boilerplates/backend/go-fiber/internal/core/task"
	pgclient "github.com/makifbaysal/boilerplates/backend/go-fiber/pkg/postgres"
)

// Run via docker-compose.dev.yml — see that file for the exact command.
// Requires POSTGRES_TEST_DSN and the schema in schema.sql applied.
func TestTaskRepository_CRUD(t *testing.T) {
	dsn := os.Getenv("POSTGRES_TEST_DSN")
	if dsn == "" {
		t.Skip("POSTGRES_TEST_DSN not set, skipping postgres integration test")
	}

	ctx := context.Background()
	pool, err := pgclient.Connect(ctx, pgclient.Config{DSN: dsn})
	require.NoError(t, err)
	defer pool.Close()

	_, err = pool.Exec(ctx, `CREATE TABLE IF NOT EXISTS tasks (
		id TEXT PRIMARY KEY, title TEXT NOT NULL, done BOOLEAN NOT NULL DEFAULT FALSE,
		created_at TIMESTAMPTZ NOT NULL, updated_at TIMESTAMPTZ NOT NULL)`)
	require.NoError(t, err)

	repo := postgres.NewTaskRepository(pool)
	id := uuid.NewString()
	t.Cleanup(func() { _, _ = pool.Exec(ctx, `DELETE FROM tasks WHERE id = $1`, id) })

	now := time.Now().UTC().Truncate(time.Microsecond)
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
