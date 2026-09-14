// Package cassandra is a reusable, app-agnostic gocql session
// constructor. Outbound adapters (internal/adapter/outbound/cassandra)
// build on top of this.
package cassandra

import (
	"fmt"
	"time"

	"github.com/gocql/gocql"
)

type Config struct {
	Hosts             []string
	Keyspace          string
	ConnectTimeout    time.Duration
	ConsistencyString string // one of gocql's consistency names, default "quorum"
}

func Connect(cfg Config) (*gocql.Session, error) {
	cluster := gocql.NewCluster(cfg.Hosts...)
	cluster.Keyspace = cfg.Keyspace

	if cfg.ConnectTimeout > 0 {
		cluster.ConnectTimeout = cfg.ConnectTimeout
	} else {
		cluster.ConnectTimeout = 5 * time.Second
	}

	consistency := cfg.ConsistencyString
	if consistency == "" {
		consistency = "quorum"
	}
	level, err := gocql.ParseConsistencyWrapper(consistency)
	if err != nil {
		return nil, fmt.Errorf("cassandra: invalid consistency %q: %w", consistency, err)
	}
	cluster.Consistency = level

	session, err := cluster.CreateSession()
	if err != nil {
		return nil, fmt.Errorf("cassandra: connect: %w", err)
	}
	return session, nil
}
