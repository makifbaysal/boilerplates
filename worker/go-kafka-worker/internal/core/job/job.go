// Package job is the hexagon core: domain model, ports, and use-case
// logic for processing async work items. Imports nothing
// framework-specific (no kafka/rabbitmq client, no db driver) — that's
// what lets kafka and rabbitmq consumers share one implementation, the
// same way rest/grpc/graphql share task.UseCase in ../../../backend/go-fiber.
package job

import "time"

type Job struct {
	ID          string
	TaskID      string
	Type        string
	Payload     string
	ReceivedAt  time.Time
	ProcessedAt time.Time
}
