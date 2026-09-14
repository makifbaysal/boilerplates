package ai.tasktrooper.boilerplate.quarkus.core.task;

import java.util.List;
import java.util.Optional;

/**
 * Driven (outbound) port: what the core needs from persistence.
 * Implemented by adapter.outbound.{memory,postgres,mongo,cassandra}.
 * Mocked via Mockito for TaskService tests.
 */
public interface TaskRepositoryPort {
    Task create(Task task);

    Optional<Task> get(String id);

    List<Task> list();

    Task update(Task task);

    void delete(String id);
}
