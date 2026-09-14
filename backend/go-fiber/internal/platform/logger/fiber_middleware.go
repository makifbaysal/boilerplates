package logger

import (
	"log/slog"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

// FiberMiddleware logs one structured line per request (method, path,
// status, latency, request_id) and attaches a request-scoped logger to
// the request context via SetContext, so downstream code reaches it
// through logger.FromContext(c.Context()).
func FiberMiddleware(log *slog.Logger) fiber.Handler {
	return func(c fiber.Ctx) error {
		start := time.Now()

		reqID := c.Get(fiber.HeaderXRequestID)
		if reqID == "" {
			reqID = uuid.NewString()
		}
		c.Set(fiber.HeaderXRequestID, reqID)

		reqLog := log.With("request_id", reqID)
		c.SetContext(WithContext(c.Context(), reqLog))

		err := c.Next()

		reqLog.Info("http_request",
			"method", c.Method(),
			"path", c.Path(),
			"status", c.Response().StatusCode(),
			"duration_ms", time.Since(start).Milliseconds(),
		)
		return err
	}
}
