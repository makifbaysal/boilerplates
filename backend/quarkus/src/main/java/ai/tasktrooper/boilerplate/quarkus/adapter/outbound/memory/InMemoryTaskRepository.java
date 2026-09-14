package ai.tasktrooper.boilerplate.quarkus.adapter.outbound.memory;

import java.util.Comparator;
import java.util.List;
import java.util.Map;
import java.util.Optional;
import java.util.concurrent.ConcurrentHashMap;

import jakarta.enterprise.context.ApplicationScoped;
import jakarta.enterprise.inject.Typed;

import ai.tasktrooper.boilerplate.quarkus.core.task.Task;
import ai.tasktrooper.boilerplate.quarkus.core.task.TaskNotFoundException;
import ai.tasktrooper.boilerplate.quarkus.core.task.TaskRepositoryPort;

/**
 * In-memory TaskRepositoryPort — the default when {@code app.datastore}
 * is unset or "memory". Swap it by setting {@code app.datastore} to
 * postgres/mongo/cassandra; nothing else changes, that's the point of
 * the port.
 *
 * <p>{@code @Typed(InMemoryTaskRepository.class)} restricts this bean's
 * CDI types to just its own concrete class — without it, every one of
 * the four TaskRepositoryPort implementations is a candidate for any
 * {@code @Inject TaskRepositoryPort} injection point simultaneously,
 * which is an AmbiguousResolutionException at build time (this actually
 * happened — {@code app.datastore} is a runtime-overridable property,
 * not a build-time one, so ArC can't use it to exclude beans the way
 * {@code @LookupIfProperty} does for build-time properties). See
 * TaskRepositorySelector, which is the one place that resolves the
 * runtime choice, and .ai/architecture.md#choosing-a-database.
 */
@ApplicationScoped
@Typed(InMemoryTaskRepository.class)
public class InMemoryTaskRepository implements TaskRepositoryPort {

    private final Map<String, Task> tasks = new ConcurrentHashMap<>();

    @Override
    public Task create(Task task) {
        tasks.put(task.id(), task);
        return task;
    }

    @Override
    public Optional<Task> get(String id) {
        return Optional.ofNullable(tasks.get(id));
    }

    @Override
    public List<Task> list() {
        return tasks.values().stream().sorted(Comparator.comparing(Task::createdAt)).toList();
    }

    @Override
    public Task update(Task task) {
        if (!tasks.containsKey(task.id())) {
            throw new TaskNotFoundException(task.id());
        }
        tasks.put(task.id(), task);
        return task;
    }

    @Override
    public void delete(String id) {
        if (tasks.remove(id) == null) {
            throw new TaskNotFoundException(id);
        }
    }
}
