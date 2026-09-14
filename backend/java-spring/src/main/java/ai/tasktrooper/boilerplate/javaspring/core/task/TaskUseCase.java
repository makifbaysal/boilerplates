package ai.tasktrooper.boilerplate.javaspring.core.task;

import java.util.List;

/**
 * Driving (inbound) port: what every inbound adapter (rest, grpc,
 * graphql) is allowed to call. Adapters depend on this interface, never
 * on {@link TaskService} directly.
 */
public interface TaskUseCase {
    Task create(String title);

    Task get(String id);

    List<Task> list();

    Task setDone(String id, boolean done);

    void delete(String id);
}
