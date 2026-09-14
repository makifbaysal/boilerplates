// Package elasticsearch is a reusable, app-agnostic go-elasticsearch
// typed-client constructor. Not wired to any Repository — use it for
// search/analytics indexes fed from task events, not as the primary
// store.
package elasticsearch

import (
	"context"
	"fmt"
	"time"

	elasticsearch "github.com/elastic/go-elasticsearch/v9"
)

type Config struct {
	Addresses      []string
	Username       string
	Password       string
	ConnectTimeout time.Duration
}

func Connect(ctx context.Context, cfg Config) (*elasticsearch.TypedClient, error) {
	opts := []elasticsearch.Option{elasticsearch.WithAddresses(cfg.Addresses...)}
	if cfg.Username != "" {
		opts = append(opts, elasticsearch.WithBasicAuth(cfg.Username, cfg.Password))
	}

	client, err := elasticsearch.NewTyped(opts...)
	if err != nil {
		return nil, fmt.Errorf("elasticsearch: new client: %w", err)
	}

	timeout := cfg.ConnectTimeout
	if timeout <= 0 {
		timeout = 5 * time.Second
	}
	pingCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	ok, err := client.Ping().Do(pingCtx)
	if err != nil {
		return nil, fmt.Errorf("elasticsearch: ping: %w", err)
	}
	if !ok {
		return nil, fmt.Errorf("elasticsearch: ping: cluster unreachable")
	}
	return client, nil
}
