// Package config loads layered, environment-aware configuration with
// koanf: configs/config.<APP_ENV>.yaml as the base, APP_-prefixed
// environment variables as overrides, with optional live file-watch
// reload. See .ai/common-tasks.md for "add a config field".
package config

type Config struct {
	Env       string          `koanf:"env"`
	Server    ServerConfig    `koanf:"server"`
	Log       LogConfig       `koanf:"log"`
	RateLimit RateLimitConfig `koanf:"rate_limit"`
	Breaker   BreakerConfig   `koanf:"breaker"`
}

// RateLimitConfig caps inbound requests per client IP. Enabled by default:
// the limiter is only skipped when a config explicitly disables it.
type RateLimitConfig struct {
	Disabled      bool `koanf:"disabled"`
	Max           int  `koanf:"max"`
	WindowSeconds int  `koanf:"window_seconds"`
}

// BreakerConfig is the default circuit breaker setting applied to outbound
// dependencies (see pkg/resilience).
type BreakerConfig struct {
	FailureThreshold int `koanf:"failure_threshold"`
	OpenSeconds      int `koanf:"open_seconds"`
	HalfOpenMaxCalls int `koanf:"half_open_max_calls"`
}

type ServerConfig struct {
	Port string `koanf:"port"`
}

type LogConfig struct {
	Level  string `koanf:"level"`
	Format string `koanf:"format"`
}
