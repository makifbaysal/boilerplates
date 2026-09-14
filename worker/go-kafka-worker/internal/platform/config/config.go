// Package config loads layered, environment-aware configuration with
// koanf: configs/config.<APP_ENV>.yaml as the base, APP_-prefixed
// environment variables as overrides, with optional live file-watch
// reload. See .ai/common-tasks.md for "add a config field".
package config

type Config struct {
	Env       string          `koanf:"env"`
	Kafka     KafkaConfig     `koanf:"kafka"`
	RabbitMQ  RabbitMQConfig  `koanf:"rabbitmq"`
	Log       LogConfig       `koanf:"log"`
	RateLimit RateLimitConfig `koanf:"rate_limit"`
	Breaker   BreakerConfig   `koanf:"breaker"`
}

// RateLimitConfig throttles how fast messages are processed. A consumer with
// no ceiling drains a backlog as fast as the broker can serve it and takes the
// downstream database down with it; the queue is a buffer, so waiting is free.
type RateLimitConfig struct {
	// PerSecond is the sustained processing rate; 0 disables throttling.
	PerSecond float64 `koanf:"per_second"`
	// Burst is how many messages may be processed back-to-back.
	Burst float64 `koanf:"burst"`
	// WaitSeconds is how long a message waits for a token before being
	// processed anyway (0 = wait is skipped).
	WaitSeconds int `koanf:"wait_seconds"`
}

// BreakerConfig is the default circuit breaker setting applied to outbound
// dependencies (see pkg/resilience).
type BreakerConfig struct {
	FailureThreshold int `koanf:"failure_threshold"`
	OpenSeconds      int `koanf:"open_seconds"`
	HalfOpenMaxCalls int `koanf:"half_open_max_calls"`
}

type KafkaConfig struct {
	Brokers []string `koanf:"brokers"`
	Topic   string   `koanf:"topic"`
	GroupID string   `koanf:"group_id"`
}

type RabbitMQConfig struct {
	URL   string `koanf:"url"`
	Queue string `koanf:"queue"`
}

type LogConfig struct {
	Level  string `koanf:"level"`
	Format string `koanf:"format"`
}
