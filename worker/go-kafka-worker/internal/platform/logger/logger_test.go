package logger_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/makifbaysal/boilerplates/worker/go-kafka-worker/internal/platform/logger"
)

func TestFromContext_DefaultsWhenUnset(t *testing.T) {
	got := logger.FromContext(context.Background())
	require.Equal(t, slog.Default(), got)
}

func TestWithLogging_LogsAndAttachesRequestScopedLogger(t *testing.T) {
	var buf bytes.Buffer
	log := slog.New(slog.NewJSONHandler(&buf, nil))

	var sawMessageID any
	handler := logger.WithLogging(log, func(ctx context.Context, source string, value []byte) error {
		sawMessageID = logger.FromContext(ctx) != slog.Default()
		return nil
	})

	err := handler(context.Background(), "test-topic", []byte(`{"ok":true}`))
	require.NoError(t, err)
	require.Equal(t, true, sawMessageID)

	var line map[string]any
	require.NoError(t, json.Unmarshal(bytes.TrimSpace(buf.Bytes()), &line))
	require.Equal(t, "message_processed", line["msg"])
	require.Equal(t, "test-topic", line["source"])
	require.NotEmpty(t, line["message_id"])
}

func TestWithLogging_LogsErrorLevelOnFailure(t *testing.T) {
	var buf bytes.Buffer
	log := slog.New(slog.NewJSONHandler(&buf, nil))

	wantErr := errors.New("boom")
	handler := logger.WithLogging(log, func(ctx context.Context, source string, value []byte) error {
		return wantErr
	})

	err := handler(context.Background(), "test-topic", nil)
	require.ErrorIs(t, err, wantErr)

	var line map[string]any
	require.NoError(t, json.Unmarshal(bytes.TrimSpace(buf.Bytes()), &line))
	require.Equal(t, "ERROR", line["level"])
	require.Equal(t, "boom", line["error"])
}
