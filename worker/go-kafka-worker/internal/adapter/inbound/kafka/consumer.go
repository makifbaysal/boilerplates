// Package kafka adapts franz-go to the job.Processor input port — the
// Kafka sibling of adapter/inbound/rabbitmq.
package kafka

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"

	"github.com/twmb/franz-go/pkg/kgo"

	"github.com/makifbaysal/boilerplates/worker/go-kafka-worker/internal/core/job"
	"github.com/makifbaysal/boilerplates/worker/go-kafka-worker/internal/platform/logger"
	"github.com/makifbaysal/boilerplates/worker/go-kafka-worker/pkg/resilience"
)

type message struct {
	TaskID string `json:"task_id"`
	Type   string `json:"type"`
}

type Consumer struct {
	client  *kgo.Client
	proc    job.Processor
	log     *slog.Logger
	handle  logger.MessageHandler
	limiter *resilience.Limiter
	wait    time.Duration
}

func NewConsumer(client *kgo.Client, proc job.Processor, log *slog.Logger) *Consumer {
	c := &Consumer{client: client, proc: proc, log: log}
	c.handle = logger.WithLogging(log, c.HandleMessage)
	return c
}

// WithRateLimit caps the processing rate. Without it the consumer runs as fast
// as the broker delivers, which is exactly how a backlog replay turns into a
// database outage.
func (c *Consumer) WithRateLimit(perSecond, burst float64, wait time.Duration) *Consumer {
	c.limiter = resilience.NewLimiter(perSecond, burst)
	c.wait = wait
	return c
}

// HandleMessage decodes value and hands it to job.Processor. Exported
// and dependency-free beyond job.Processor, so it's directly
// unit-testable without a live broker — see consumer_test.go. Wrapped
// with logger.WithLogging in NewConsumer for the structured per-message
// log line; call this directly only from a test.
func (c *Consumer) HandleMessage(ctx context.Context, _ string, value []byte) error {
	// Throttle before doing any work: a queued job can wait, a starved
	// database cannot.
	if c.limiter != nil && c.wait > 0 && !c.limiter.Wait(c.wait) {
		c.log.Warn("rate limit wait expired, processing anyway")
	}
	var msg message
	if err := json.Unmarshal(value, &msg); err != nil {
		return err
	}
	_, err := c.proc.Process(ctx, msg.TaskID, msg.Type, string(value))
	return err
}

// Run polls until ctx is cancelled. A handler error is logged (inside
// WithLogging) and the offset is still committed — this boilerplate
// treats a bad/unprocessable message as skip-and-move-on rather than
// retry-forever or dead-letter; swap that policy in HandleMessage's
// caller if your use case needs at-least-once + DLQ semantics.
func (c *Consumer) Run(ctx context.Context) error {
	for {
		if ctx.Err() != nil {
			return ctx.Err()
		}

		fetches := c.client.PollFetches(ctx)
		if ctx.Err() != nil {
			return ctx.Err()
		}

		fetches.EachError(func(topic string, partition int32, err error) {
			c.log.Error("fetch error", "topic", topic, "partition", partition, "error", err)
		})

		fetches.EachRecord(func(r *kgo.Record) {
			_ = c.handle(ctx, r.Topic, r.Value)
		})

		if err := c.client.CommitUncommittedOffsets(ctx); err != nil {
			c.log.Error("commit failed", "error", err)
		}
	}
}
