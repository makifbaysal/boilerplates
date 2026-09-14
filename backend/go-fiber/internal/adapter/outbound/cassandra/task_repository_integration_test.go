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

	"github.com/makifbaysal/boilerplates/backend/go-fiber/internal/adapter/outbound/cassandra"
	"github.com/makifbaysal/boilerplates/backend/go-fiber/internal/core/task"
	cassandraclient "github.com/makifbaysal/boilerplates/backend/go-fiber/pkg/cassandra"
)

// Run via docker-compose.dev.yml — see that file for the exact command.
// Requires the keyspace/table in schema.cql to already exist.
func TestTaskRepository_CRUD(t *testing.T) {
	hosts := os.Getenv("CASSANDRA_TEST_HOSTS")
	if hosts == "" {
		t.Skip("CASSANDRA_TEST_HOSTS not set, skipping cassandra integration test")
	}

	session, err := cassandraclient.Connect(cassandraclient.Config{
		Hosts:    strings.Split(hosts, ","),
		Keyspace: "boilerplate_test",
	})
	require.NoError(t, err)
	defer session.Close()

	require.NoError(t, session.Query(`CREATE TABLE IF NOT EXISTS tasks (
		id text PRIMARY KEY, title text, done boolean, created_at timestamp, updated_at timestamp)`).Exec())

	repo := cassandra.NewTaskRepository(session)
	ctx := context.Background()
	id := uuid.NewString()
	t.Cleanup(func() { _ = session.Query(`DELETE FROM tasks WHERE id = ?`, id).Exec() })

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

	require.NoError(t, repo.Delete(ctx, id))
	_, err = repo.Get(ctx, id)
	require.ErrorIs(t, err, task.ErrNotFound)
}
