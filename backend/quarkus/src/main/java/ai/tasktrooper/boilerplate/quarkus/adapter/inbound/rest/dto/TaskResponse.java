package ai.tasktrooper.boilerplate.quarkus.adapter.inbound.rest.dto;

import java.time.Instant;

import ai.tasktrooper.boilerplate.quarkus.core.task.Task;

public record TaskResponse(String id, String title, boolean done, Instant createdAt, Instant updatedAt) {
    public static TaskResponse from(Task task) {
        return new TaskResponse(task.id(), task.title(), task.done(), task.createdAt(), task.updatedAt());
    }
}
