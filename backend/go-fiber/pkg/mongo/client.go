// Package mongo is a reusable, app-agnostic mongo-driver/v2 client
// constructor. Outbound adapters (internal/adapter/outbound/mongo) build
// on top of this.
package mongo

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type Config struct {
	URI            string // mongodb://user:pass@host:27017
	Database       string
	ConnectTimeout time.Duration
}

func Connect(ctx context.Context, cfg Config) (*mongo.Client, error) {
	timeout := cfg.ConnectTimeout
	if timeout <= 0 {
		timeout = 5 * time.Second
	}
	connectCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	client, err := mongo.Connect(options.Client().ApplyURI(cfg.URI))
	if err != nil {
		return nil, fmt.Errorf("mongo: connect: %w", err)
	}
	if err := client.Ping(connectCtx, nil); err != nil {
		_ = client.Disconnect(connectCtx)
		return nil, fmt.Errorf("mongo: ping: %w", err)
	}
	return client, nil
}
