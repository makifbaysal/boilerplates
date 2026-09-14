//go:build integration

package rabbitmq_test

import (
	"context"
	"log/slog"
	"os"
	"testing"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/stretchr/testify/require"

	rabbitmqadapter "github.com/makifbaysal/boilerplates/worker/go-kafka-worker/internal/adapter/inbound/rabbitmq"
	"github.com/makifbaysal/boilerplates/worker/go-kafka-worker/internal/adapter/outbound/memory"
	"github.com/makifbaysal/boilerplates/worker/go-kafka-worker/internal/core/job"
)

// Run via docker-compose.dev.yml — see that file for the exact command.
// Proves the channel/delivery wiring against a real broker;
// HandleMessage's decode/process logic is already covered by
// consumer_test.go without needing a broker at all.
func TestConsumer_Run(t *testing.T) {
	url := os.Getenv("RABBITMQ_TEST_URL")
	if url == "" {
		t.Skip("RABBITMQ_TEST_URL not set, skipping rabbitmq integration test")
	}

	conn, err := amqp.Dial(url)
	require.NoError(t, err)
	defer conn.Close()

	ch, err := conn.Channel()
	require.NoError(t, err)
	defer ch.Close()

	queue := "worker-integration-test"
	_, err = ch.QueueDeclare(queue, false, true, false, false, nil)
	require.NoError(t, err)

	require.NoError(t, ch.Publish("", queue, false, false, amqp.Publishing{
		Body: []byte(`{"task_id":"task-1","type":"task.created"}`),
	}))

	deliveries, err := ch.Consume(queue, "worker-integration-test-consumer", false, false, false, false, nil)
	require.NoError(t, err)

	store := memory.NewJobStore()
	uc := job.NewService(store)
	consumer := rabbitmqadapter.NewConsumer(uc, slog.Default())

	runCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	done := make(chan struct{})
	go func() {
		_ = consumer.Run(runCtx, queue, deliveries)
		close(done)
	}()

	require.Eventually(t, func() bool {
		list, _ := store.List(context.Background())
		return len(list) == 1
	}, 9*time.Second, 200*time.Millisecond)

	cancel()
	<-done
}
