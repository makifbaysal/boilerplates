package ai.tasktrooper.boilerplate.javaspring.adapter.outbound.mongo;

import java.util.List;
import java.util.Optional;

import org.springframework.context.annotation.Profile;
import org.springframework.stereotype.Repository;

import ai.tasktrooper.boilerplate.javaspring.core.task.Task;
import ai.tasktrooper.boilerplate.javaspring.core.task.TaskNotFoundException;
import ai.tasktrooper.boilerplate.javaspring.core.task.TaskRepositoryPort;

/** Activate with {@code SPRING_PROFILES_ACTIVE=mongo}. */
@Repository
@Profile("mongo")
public class MongoTaskRepository implements TaskRepositoryPort {

    private final TaskMongoRepository mongoRepository;

    public MongoTaskRepository(TaskMongoRepository mongoRepository) {
        this.mongoRepository = mongoRepository;
    }

    private static Task toDomain(TaskMongoDocument d) {
        return new Task(d.getId(), d.getTitle(), d.isDone(), d.getCreatedAt(), d.getUpdatedAt());
    }

    private static TaskMongoDocument toDocument(Task t) {
        return new TaskMongoDocument(t.id(), t.title(), t.done(), t.createdAt(), t.updatedAt());
    }

    @Override
    public Task create(Task task) {
        return toDomain(mongoRepository.save(toDocument(task)));
    }

    @Override
    public Optional<Task> get(String id) {
        return mongoRepository.findById(id).map(MongoTaskRepository::toDomain);
    }

    @Override
    public List<Task> list() {
        return mongoRepository.findAll().stream().map(MongoTaskRepository::toDomain).toList();
    }

    @Override
    public Task update(Task task) {
        if (!mongoRepository.existsById(task.id())) {
            throw new TaskNotFoundException(task.id());
        }
        return toDomain(mongoRepository.save(toDocument(task)));
    }

    @Override
    public void delete(String id) {
        if (!mongoRepository.existsById(id)) {
            throw new TaskNotFoundException(id);
        }
        mongoRepository.deleteById(id);
    }
}
