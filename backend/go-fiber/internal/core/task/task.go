// Package task is the hexagon core: domain model, ports, and use-case
// logic. It imports nothing framework-specific (no fiber, no grpc, no
// gqlgen, no DB driver) — that's what keeps it testable and reusable
// across every inbound/outbound adapter.
package task

import "time"

type Task struct {
	ID        string
	Title     string
	Done      bool
	CreatedAt time.Time
	UpdatedAt time.Time
}
