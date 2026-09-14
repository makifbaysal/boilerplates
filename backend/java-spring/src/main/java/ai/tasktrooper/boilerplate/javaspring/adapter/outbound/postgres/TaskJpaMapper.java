package ai.tasktrooper.boilerplate.javaspring.adapter.outbound.postgres;

import org.mapstruct.Mapper;

import ai.tasktrooper.boilerplate.javaspring.core.task.Task;

/**
 * MapStruct generates the implementation at compile time (see
 * target/generated-sources/annotations after `mvn compile`) — the
 * code-gen counterpart to Lombok on the entities. mongo/cassandra
 * adapters map by hand since they're a single trivial method each; this
 * one demonstrates the pattern for a resource with more fields/adapters
 * to keep in sync.
 */
@Mapper(componentModel = "spring")
public interface TaskJpaMapper {
    Task toDomain(TaskJpaEntity entity);

    TaskJpaEntity toEntity(Task task);
}
