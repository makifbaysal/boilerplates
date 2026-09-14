package ai.tasktrooper.boilerplate.javaspring.core.task;

import java.time.Instant;
import java.util.List;
import java.util.UUID;

import org.springframework.stereotype.Service;

@Service
public class TaskService implements TaskUseCase {

    private final TaskRepositoryPort repository;

    public TaskService(TaskRepositoryPort repository) {
        this.repository = repository;
    }

    @Override
    public Task create(String title) {
        String trimmed = title == null ? "" : title.trim();
        if (trimmed.isEmpty()) {
            throw new InvalidTaskTitleException();
        }
        Instant now = Instant.now();
        Task task = new Task(UUID.randomUUID().toString(), trimmed, false, now, now);
        return repository.create(task);
    }

    @Override
    public Task get(String id) {
        return repository.get(id).orElseThrow(() -> new TaskNotFoundException(id));
    }

    @Override
    public List<Task> list() {
        return repository.list();
    }

    @Override
    public Task setDone(String id, boolean done) {
        Task existing = get(id);
        Task updated = existing.withDone(done, Instant.now());
        return repository.update(updated);
    }

    @Override
    public void delete(String id) {
        repository.delete(id);
    }
}
