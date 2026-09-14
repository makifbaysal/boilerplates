package ai.tasktrooper.boilerplate.quarkus.core.task;

import java.time.Instant;

/**
 * Domain model. No JPA/Jackson/protobuf annotations here — adapters own
 * their own DTOs/entities and map to/from this, the same rule as
 * backend/go-fiber's core.Task.
 */
public record Task(String id, String title, boolean done, Instant createdAt, Instant updatedAt) {

    public Task withDone(boolean newDone, Instant updatedAt) {
        return new Task(id, title, newDone, createdAt, updatedAt);
    }
}
