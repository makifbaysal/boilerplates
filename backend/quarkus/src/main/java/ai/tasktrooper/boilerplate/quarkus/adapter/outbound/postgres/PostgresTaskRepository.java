package ai.tasktrooper.boilerplate.quarkus.adapter.outbound.postgres;

import java.util.List;
import java.util.Optional;

import jakarta.enterprise.context.ApplicationScoped;
import jakarta.enterprise.inject.Typed;
import jakarta.inject.Inject;
import jakarta.transaction.Transactional;

import ai.tasktrooper.boilerplate.quarkus.core.task.Task;
import ai.tasktrooper.boilerplate.quarkus.core.task.TaskNotFoundException;
import ai.tasktrooper.boilerplate.quarkus.core.task.TaskRepositoryPort;

/** Activate with {@code app.datastore=postgres} — see
 * TaskRepositorySelector and memory's InMemoryTaskRepository javadoc for
 * why {@code @Typed} is here instead of a property-conditional
 * annotation. */
@ApplicationScoped
@Typed(PostgresTaskRepository.class)
public class PostgresTaskRepository implements TaskRepositoryPort {

    private final TaskPanacheRepository repository;

    @Inject
    public PostgresTaskRepository(TaskPanacheRepository repository) {
        this.repository = repository;
    }

    private static Task toDomain(TaskPanacheEntity e) {
        return new Task(e.id, e.title, e.done, e.createdAt, e.updatedAt);
    }

    private static TaskPanacheEntity toEntity(Task t) {
        TaskPanacheEntity e = new TaskPanacheEntity();
        e.id = t.id();
        e.title = t.title();
        e.done = t.done();
        e.createdAt = t.createdAt();
        e.updatedAt = t.updatedAt();
        return e;
    }

    @Override
    @Transactional
    public Task create(Task task) {
        TaskPanacheEntity e = toEntity(task);
        repository.persist(e);
        return toDomain(e);
    }

    @Override
    public Optional<Task> get(String id) {
        return repository.findByIdOptional(id).map(PostgresTaskRepository::toDomain);
    }

    @Override
    public List<Task> list() {
        return repository.listAll().stream().map(PostgresTaskRepository::toDomain).toList();
    }

    @Override
    @Transactional
    public Task update(Task task) {
        TaskPanacheEntity e = repository.findByIdOptional(task.id())
                .orElseThrow(() -> new TaskNotFoundException(task.id()));
        e.title = task.title();
        e.done = task.done();
        e.updatedAt = task.updatedAt();
        return toDomain(e);
    }

    @Override
    @Transactional
    public void delete(String id) {
        boolean deleted = repository.deleteById(id);
        if (!deleted) {
            throw new TaskNotFoundException(id);
        }
    }
}
