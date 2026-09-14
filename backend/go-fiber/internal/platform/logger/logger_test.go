package logger_test

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v3"
	"github.com/stretchr/testify/require"

	"github.com/makifbaysal/boilerplates/backend/go-fiber/internal/platform/logger"
)

func TestFromContext_DefaultsWhenUnset(t *testing.T) {
	got := logger.FromContext(context.Background())
	require.Equal(t, slog.Default(), got)
}

func TestFiberMiddleware_LogsAndSetsRequestID(t *testing.T) {
	var buf bytes.Buffer
	log := slog.New(slog.NewJSONHandler(&buf, nil))

	app := fiber.New()
	app.Use(logger.FiberMiddleware(log))
	app.Get("/ping", func(c fiber.Ctx) error { return c.SendString("pong") })

	resp, err := app.Test(httptest.NewRequest(http.MethodGet, "/ping", nil))
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	require.NotEmpty(t, resp.Header.Get(fiber.HeaderXRequestID))

	var line map[string]any
	require.NoError(t, json.Unmarshal(bytes.TrimSpace(buf.Bytes()), &line))
	require.Equal(t, "http_request", line["msg"])
	require.Equal(t, float64(http.StatusOK), line["status"])
	require.NotEmpty(t, line["request_id"])
}
