package ai.tasktrooper.boilerplate.quarkus.core.task;

/** Sentinel domain exception for a blank/missing title. */
public class InvalidTaskTitleException extends RuntimeException {
    public InvalidTaskTitleException() {
        super("task title is required");
    }
}
