package ai.tasktrooper.boilerplate.javaspring.core.task;

/** Sentinel domain exception, checked by type — the Java equivalent of
 * go-fiber's {@code errors.Is(err, task.ErrNotFound)}. */
public class TaskNotFoundException extends RuntimeException {
    public TaskNotFoundException(String id) {
        super("task not found: " + id);
    }
}
