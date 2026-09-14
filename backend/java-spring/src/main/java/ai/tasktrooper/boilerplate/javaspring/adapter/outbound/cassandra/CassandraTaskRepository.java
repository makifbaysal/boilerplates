package ai.tasktrooper.boilerplate.javaspring.adapter.outbound.cassandra;

import java.util.List;
import java.util.Optional;

import org.springframework.context.annotation.Profile;
import org.springframework.stereotype.Repository;

import ai.tasktrooper.boilerplate.javaspring.core.task.Task;
import ai.tasktrooper.boilerplate.javaspring.core.task.TaskNotFoundException;
import ai.tasktrooper.boilerplate.javaspring.core.task.TaskRepositoryPort;

/** Activate with {@code SPRING_PROFILES_ACTIVE=cassandra}. */
@Repository
@Profile("cassandra")
public class CassandraTaskRepository implements TaskRepositoryPort {

    private final TaskCassandraRepository cassandraRepository;

    public CassandraTaskRepository(TaskCassandraRepository cassandraRepository) {
        this.cassandraRepository = cassandraRepository;
    }

    private static Task toDomain(TaskCassandraEntity e) {
        return new Task(e.getId(), e.getTitle(), e.isDone(), e.getCreatedAt(), e.getUpdatedAt());
    }

    private static TaskCassandraEntity toEntity(Task t) {
        return new TaskCassandraEntity(t.id(), t.title(), t.done(), t.createdAt(), t.updatedAt());
    }

    @Override
    public Task create(Task task) {
        return toDomain(cassandraRepository.save(toEntity(task)));
    }

    @Override
    public Optional<Task> get(String id) {
        return cassandraRepository.findById(id).map(CassandraTaskRepository::toDomain);
    }

    @Override
    public List<Task> list() {
        return cassandraRepository.findAll().stream().map(CassandraTaskRepository::toDomain).toList();
    }

    // Cassandra's UPDATE is an upsert at the CQL level — check existence
    // first so this port's "update must already exist" contract holds
    // the same way it does for postgres/mongo.
    @Override
    public Task update(Task task) {
        if (!cassandraRepository.existsById(task.id())) {
            throw new TaskNotFoundException(task.id());
        }
        return toDomain(cassandraRepository.save(toEntity(task)));
    }

    @Override
    public void delete(String id) {
        if (!cassandraRepository.existsById(id)) {
            throw new TaskNotFoundException(id);
        }
        cassandraRepository.deleteById(id);
    }
}
