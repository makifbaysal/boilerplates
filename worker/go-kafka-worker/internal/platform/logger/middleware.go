package logger

import (
	"context"
	"log/slog"
	"time"

	"github.com/google/uuid"
)

// MessageHandler is the shape every queue consumer's per-message handler
// implements — kafka.Consumer.HandleMessage, rabbitmq.Consumer.HandleMessage.
// source is the topic (kafka) or queue name (rabbitmq).
type MessageHandler func(ctx context.Context, source string, value []byte) error

// WithLogging wraps a MessageHandler with one structured log line per
// message (source, size, duration, error) and attaches a request-scoped
// logger (with message_id) to ctx — the queue-consumer equivalent of
// rest.FiberMiddleware / grpc's UnaryServerInterceptor in the go-fiber
// boilerplate.
func WithLogging(log *slog.Logger, next MessageHandler) MessageHandler {
	return func(ctx context.Context, source string, value []byte) error {
		start := time.Now()
		msgLog := log.With("message_id", uuid.NewString())
		ctx = WithContext(ctx, msgLog)

		err := next(ctx, source, value)

		attrs := []any{
			"source", source,
			"bytes", len(value),
			"duration_ms", time.Since(start).Milliseconds(),
		}
		if err != nil {
			msgLog.Error("message_processed", append(attrs, "error", err.Error())...)
		} else {
			msgLog.Info("message_processed", attrs...)
		}
		return err
	}
}
