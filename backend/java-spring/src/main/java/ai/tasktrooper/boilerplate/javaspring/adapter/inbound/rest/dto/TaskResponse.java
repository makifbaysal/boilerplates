package ai.tasktrooper.boilerplate.javaspring.adapter.inbound.rest.dto;

import java.time.Instant;

import ai.tasktrooper.boilerplate.javaspring.core.task.Task;

/** This adapter's own wire shape — never expose core.Task directly. */
public record TaskResponse(String id, String title, boolean done, Instant createdAt, Instant updatedAt) {
    public static TaskResponse from(Task task) {
        return new TaskResponse(task.id(), task.title(), task.done(), task.createdAt(), task.updatedAt());
    }
}
