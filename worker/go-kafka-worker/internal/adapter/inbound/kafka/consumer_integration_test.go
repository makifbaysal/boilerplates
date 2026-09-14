//go:build integration

package kafka_test

import (
	"context"
	"log/slog"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/twmb/franz-go/pkg/kgo"

	kafkaadapter "github.com/makifbaysal/boilerplates/worker/go-kafka-worker/internal/adapter/inbound/kafka"
	"github.com/makifbaysal/boilerplates/worker/go-kafka-worker/internal/adapter/outbound/memory"
	"github.com/makifbaysal/boilerplates/worker/go-kafka-worker/internal/core/job"
)

// Run via docker-compose.dev.yml — see that file for the exact command.
// Proves the PollFetches/consumer-group wiring against a real broker;
// HandleMessage's decode/process logic is already covered by
// consumer_test.go without needing a broker at all.
func TestConsumer_Run(t *testing.T) {
	brokers := os.Getenv("KAFKA_TEST_BROKERS")
	if brokers == "" {
		t.Skip("KAFKA_TEST_BROKERS not set, skipping kafka integration test")
	}

	topic := "worker-integration-test"
	seeds := strings.Split(brokers, ",")

	producer, err := kgo.NewClient(kgo.SeedBrokers(seeds...))
	require.NoError(t, err)
	defer producer.Close()

	produceCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	res := producer.ProduceSync(produceCtx, &kgo.Record{Topic: topic, Value: []byte(`{"task_id":"task-1","type":"task.created"}`)})
	require.NoError(t, res.FirstErr())

	consumerClient, err := kgo.NewClient(
		kgo.SeedBrokers(seeds...),
		kgo.ConsumeTopics(topic),
		kgo.ConsumerGroup("worker-integration-test-group"),
		kgo.DisableAutoCommit(),
	)
	require.NoError(t, err)
	defer consumerClient.Close()

	store := memory.NewJobStore()
	uc := job.NewService(store)
	consumer := kafkaadapter.NewConsumer(consumerClient, uc, slog.Default())

	runCtx, cancelRun := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelRun()

	done := make(chan struct{})
	go func() {
		_ = consumer.Run(runCtx)
		close(done)
	}()

	require.Eventually(t, func() bool {
		list, _ := store.List(context.Background())
		return len(list) == 1
	}, 9*time.Second, 200*time.Millisecond)

	cancelRun()
	<-done
}
