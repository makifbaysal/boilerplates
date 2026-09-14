package config

import (
	"fmt"
	"log/slog"
	"os"
	"strings"
	"sync/atomic"

	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

// Loader resolves APP_ENV (default "dev") to configs/config.<env>.yaml,
// layers APP_-prefixed env vars on top as overrides, and can watch the
// file for live reload. It is safe for concurrent use.
type Loader struct {
	path string
	cur  atomic.Pointer[Config]
}

// New loads configs/config.<APP_ENV>.yaml from dir (e.g. "configs").
func New(dir string) (*Loader, error) {
	appEnv := os.Getenv("APP_ENV")
	if appEnv == "" {
		appEnv = "dev"
	}

	l := &Loader{path: fmt.Sprintf("%s/config.%s.yaml", dir, appEnv)}
	cfg, err := l.load()
	if err != nil {
		return nil, err
	}
	l.cur.Store(cfg)
	return l, nil
}

func (l *Loader) load() (*Config, error) {
	k := koanf.New(".")
	if err := k.Load(file.Provider(l.path), yaml.Parser()); err != nil {
		return nil, fmt.Errorf("config: load %s: %w", l.path, err)
	}
	if err := k.Load(envOverrides(), nil); err != nil {
		return nil, fmt.Errorf("config: load env overrides: %w", err)
	}

	var cfg Config
	if err := k.Unmarshal("", &cfg); err != nil {
		return nil, fmt.Errorf("config: unmarshal: %w", err)
	}
	return &cfg, nil
}

// envOverrides maps APP_SERVER_PORT -> server.port, APP_LOG_LEVEL ->
// log.level, etc. Env always wins over the file — it's the override
// layer, not a replacement for the demo configs.
func envOverrides() *env.Env {
	return env.Provider("APP_", ".", func(s string) string {
		return strings.ReplaceAll(strings.ToLower(strings.TrimPrefix(s, "APP_")), "_", ".")
	})
}

// Current returns the most recently loaded config. Call this per use
// rather than caching the pointer's fields, so a Watch-triggered reload
// takes effect for the next caller.
func (l *Loader) Current() *Config {
	return l.cur.Load()
}

// Watch reloads whenever the config file changes and invokes onReload
// with the new value. A failed reload is logged and the previous config
// stays in place — a bad edit never takes the process down.
func (l *Loader) Watch(log *slog.Logger, onReload func(*Config)) error {
	return file.Provider(l.path).Watch(func(_ any, err error) {
		if err != nil {
			log.Error("config watch error", "error", err)
			return
		}
		cfg, loadErr := l.load()
		if loadErr != nil {
			log.Error("config reload failed, keeping previous config", "error", loadErr)
			return
		}
		l.cur.Store(cfg)
		log.Info("config reloaded", "path", l.path)
		if onReload != nil {
			onReload(cfg)
		}
	})
}
