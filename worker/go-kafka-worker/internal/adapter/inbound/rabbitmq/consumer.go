// Package rabbitmq adapts amqp091-go to the job.Processor input port —
// the RabbitMQ sibling of adapter/inbound/kafka.
package rabbitmq

import (
	"context"
	"encoding/json"
	"log/slog"

	amqp "github.com/rabbitmq/amqp091-go"

	"github.com/makifbaysal/boilerplates/worker/go-kafka-worker/internal/core/job"
	"github.com/makifbaysal/boilerplates/worker/go-kafka-worker/internal/platform/logger"
)

type message struct {
	TaskID string `json:"task_id"`
	Type   string `json:"type"`
}

type Consumer struct {
	proc   job.Processor
	log    *slog.Logger
	handle logger.MessageHandler
}

func NewConsumer(proc job.Processor, log *slog.Logger) *Consumer {
	c := &Consumer{proc: proc, log: log}
	c.handle = logger.WithLogging(log, c.HandleMessage)
	return c
}

// HandleMessage decodes value and hands it to job.Processor. Exported
// and dependency-free beyond job.Processor, so it's directly
// unit-testable without a live broker — see consumer_test.go.
func (c *Consumer) HandleMessage(ctx context.Context, _ string, value []byte) error {
	var msg message
	if err := json.Unmarshal(value, &msg); err != nil {
		return err
	}
	_, err := c.proc.Process(ctx, msg.TaskID, msg.Type, string(value))
	return err
}

// Run ranges over deliveries until the channel closes or ctx is
// cancelled. A handler error Nacks without requeue (skip-and-move-on,
// same policy as the kafka adapter — see its Run doc for why); success
// Acks. See .ai/common-tasks.md for adding a new message type.
func (c *Consumer) Run(ctx context.Context, queue string, deliveries <-chan amqp.Delivery) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case d, ok := <-deliveries:
			if !ok {
				return nil
			}
			if err := c.handle(ctx, queue, d.Body); err != nil {
				_ = d.Nack(false, false)
				continue
			}
			_ = d.Ack(false)
		}
	}
}
