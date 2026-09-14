package ai.tasktrooper.boilerplate.quarkus.adapter.outbound.mongo;

import java.util.List;
import java.util.Optional;

import jakarta.enterprise.context.ApplicationScoped;
import jakarta.enterprise.inject.Typed;
import jakarta.inject.Inject;

import ai.tasktrooper.boilerplate.quarkus.core.task.Task;
import ai.tasktrooper.boilerplate.quarkus.core.task.TaskNotFoundException;
import ai.tasktrooper.boilerplate.quarkus.core.task.TaskRepositoryPort;

/** Activate with {@code app.datastore=mongo} — see
 * TaskRepositorySelector and memory's InMemoryTaskRepository javadoc for
 * why {@code @Typed} is here instead of a property-conditional
 * annotation. */
@ApplicationScoped
@Typed(MongoTaskRepository.class)
public class MongoTaskRepository implements TaskRepositoryPort {

    private final TaskMongoRepository repository;

    @Inject
    public MongoTaskRepository(TaskMongoRepository repository) {
        this.repository = repository;
    }

    private static Task toDomain(TaskMongoEntity e) {
        return new Task(e.id, e.title, e.done, e.createdAt, e.updatedAt);
    }

    private static TaskMongoEntity toEntity(Task t) {
        TaskMongoEntity e = new TaskMongoEntity();
        e.id = t.id();
        e.title = t.title();
        e.done = t.done();
        e.createdAt = t.createdAt();
        e.updatedAt = t.updatedAt();
        return e;
    }

    @Override
    public Task create(Task task) {
        TaskMongoEntity e = toEntity(task);
        repository.persist(e);
        return toDomain(e);
    }

    @Override
    public Optional<Task> get(String id) {
        return repository.findByIdOptional(id).map(MongoTaskRepository::toDomain);
    }

    @Override
    public List<Task> list() {
        return repository.listAll().stream().map(MongoTaskRepository::toDomain).toList();
    }

    @Override
    public Task update(Task task) {
        TaskMongoEntity e = repository.findByIdOptional(task.id())
                .orElseThrow(() -> new TaskNotFoundException(task.id()));
        e.title = task.title();
        e.done = task.done();
        e.updatedAt = task.updatedAt();
        repository.update(e);
        return toDomain(e);
    }

    @Override
    public void delete(String id) {
        boolean deleted = repository.deleteById(id);
        if (!deleted) {
            throw new TaskNotFoundException(id);
        }
    }
}
