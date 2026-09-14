// Package logger provides one structured logger (log/slog, stdlib — no
// extra dependency) shared by every inbound adapter, wired in as
// middleware/interceptor rather than passed down through the core.
package logger

import (
	"context"
	"log/slog"
	"os"
)

type ctxKey struct{}

type Config struct {
	Level  string // debug|info|warn|error, default info
	Format string // json|text, default json
}

func New(cfg Config) *slog.Logger {
	opts := &slog.HandlerOptions{Level: ParseLevel(cfg.Level)}
	return slog.New(newHandler(cfg.Format, opts))
}

// NewDynamic is New plus a *slog.LevelVar the caller can mutate at
// runtime (e.g. from config.Loader.Watch) to change the log level
// without restarting the process — see cmd/api/main.go.
func NewDynamic(cfg Config) (*slog.Logger, *slog.LevelVar) {
	var lvl slog.LevelVar
	lvl.Set(ParseLevel(cfg.Level))

	handler := newHandler(cfg.Format, &slog.HandlerOptions{Level: &lvl})
	return slog.New(handler), &lvl
}

func newHandler(format string, opts *slog.HandlerOptions) slog.Handler {
	if format == "text" {
		return slog.NewTextHandler(os.Stdout, opts)
	}
	return slog.NewJSONHandler(os.Stdout, opts)
}

// ParseLevel maps the config string to a slog.Level, defaulting to Info
// for an empty or unrecognized value.
func ParseLevel(s string) slog.Level {
	switch s {
	case "debug":
		return slog.LevelDebug
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

// WithContext attaches log to ctx so core/use-case code can pull a
// request-scoped logger (with request_id already attached) via
// FromContext instead of taking *slog.Logger as a constructor dependency
// — that would leak a transport concern into the hexagon core.
func WithContext(ctx context.Context, log *slog.Logger) context.Context {
	return context.WithValue(ctx, ctxKey{}, log)
}

// FromContext returns the request-scoped logger, or slog.Default() if
// none was attached — e.g. in a unit test that calls a use case directly
// without going through an inbound adapter's middleware/interceptor.
func FromContext(ctx context.Context) *slog.Logger {
	if log, ok := ctx.Value(ctxKey{}).(*slog.Logger); ok {
		return log
	}
	return slog.Default()
}
