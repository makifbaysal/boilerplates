package config_test

import (
	"log/slog"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/makifbaysal/boilerplates/backend/go-fiber/internal/platform/config"
)

func TestNew_ReadsDemoConfigs(t *testing.T) {
	cases := []struct {
		appEnv    string
		wantLevel string
	}{
		{appEnv: "dev", wantLevel: "debug"},
		{appEnv: "staging", wantLevel: "info"},
		{appEnv: "prod", wantLevel: "warn"},
	}

	for _, tc := range cases {
		t.Run(tc.appEnv, func(t *testing.T) {
			t.Setenv("APP_ENV", tc.appEnv)

			l, err := config.New("../../../configs")
			require.NoError(t, err)

			cfg := l.Current()
			require.Equal(t, tc.appEnv, cfg.Env)
			require.Equal(t, tc.wantLevel, cfg.Log.Level)
			require.Equal(t, "8080", cfg.Server.Port)
		})
	}
}

func TestNew_EnvOverridesFile(t *testing.T) {
	t.Setenv("APP_ENV", "dev")
	t.Setenv("APP_SERVER_PORT", "9999")

	l, err := config.New("../../../configs")
	require.NoError(t, err)

	require.Equal(t, "9999", l.Current().Server.Port)
}

func TestWatch_ReloadsOnFileChange(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.dev.yaml")
	require.NoError(t, os.WriteFile(path, []byte("env: dev\nserver:\n  port: \"8080\"\nlog:\n  level: info\n  format: json\n"), 0o644))
	t.Setenv("APP_ENV", "dev")

	l, err := config.New(dir)
	require.NoError(t, err)
	require.Equal(t, "info", l.Current().Log.Level)

	reloaded := make(chan *config.Config, 1)
	require.NoError(t, l.Watch(slog.Default(), func(c *config.Config) { reloaded <- c }))

	require.NoError(t, os.WriteFile(path, []byte("env: dev\nserver:\n  port: \"8080\"\nlog:\n  level: debug\n  format: json\n"), 0o644))

	select {
	case c := <-reloaded:
		require.Equal(t, "debug", c.Log.Level)
	case <-time.After(5 * time.Second):
		t.Fatal("timed out waiting for config reload")
	}
}
