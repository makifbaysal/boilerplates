package ai.tasktrooper.boilerplate.javaspring.adapter.inbound.graphql.dto;

import ai.tasktrooper.boilerplate.javaspring.core.task.Task;

/**
 * This adapter's own wire shape. Timestamps are ISO-8601 strings
 * (schema type String) rather than a custom scalar — a deliberate
 * simplification to avoid an extra graphql-java-extended-scalars
 * dependency for a boilerplate; swap in a real DateTime scalar if your
 * project needs typed timestamps in generated clients.
 */
public record TaskDto(String id, String title, boolean done, String createdAt, String updatedAt) {
    public static TaskDto from(Task task) {
        return new TaskDto(task.id(), task.title(), task.done(), task.createdAt().toString(), task.updatedAt().toString());
    }
}
