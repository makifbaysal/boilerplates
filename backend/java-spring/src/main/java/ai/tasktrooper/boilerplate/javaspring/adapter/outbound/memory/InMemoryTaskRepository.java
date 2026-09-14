package ai.tasktrooper.boilerplate.javaspring.adapter.outbound.memory;

import java.util.Comparator;
import java.util.List;
import java.util.Map;
import java.util.Optional;
import java.util.concurrent.ConcurrentHashMap;

import org.springframework.context.annotation.Profile;
import org.springframework.stereotype.Repository;

import ai.tasktrooper.boilerplate.javaspring.core.task.Task;
import ai.tasktrooper.boilerplate.javaspring.core.task.TaskNotFoundException;
import ai.tasktrooper.boilerplate.javaspring.core.task.TaskRepositoryPort;

/**
 * In-memory TaskRepositoryPort — the default when no {@code postgres},
 * {@code mongo}, or {@code cassandra} profile is active. Swap it for a
 * real store by activating one of those profiles; nothing else changes,
 * that's the point of the port.
 */
@Repository
@Profile("!postgres & !mongo & !cassandra")
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
