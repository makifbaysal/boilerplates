package ai.tasktrooper.boilerplate.quarkus.adapter.inbound.graphql.dto;

import java.time.Instant;

import ai.tasktrooper.boilerplate.quarkus.core.task.Task;

/** This adapter's own wire shape — SmallRye GraphQL maps java.time.Instant
 * to a DateTime scalar automatically, no custom scalar needed here
 * (unlike java-spring's ISO-string simplification). */
public record TaskDto(String id, String title, boolean done, Instant createdAt, Instant updatedAt) {
    public static TaskDto from(Task task) {
        return new TaskDto(task.id(), task.title(), task.done(), task.createdAt(), task.updatedAt());
    }
}
