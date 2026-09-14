package ai.tasktrooper.boilerplate.javaspring.adapter.outbound.postgres;

import java.util.List;
import java.util.Optional;

import org.springframework.context.annotation.Profile;
import org.springframework.stereotype.Repository;

import ai.tasktrooper.boilerplate.javaspring.core.task.Task;
import ai.tasktrooper.boilerplate.javaspring.core.task.TaskNotFoundException;
import ai.tasktrooper.boilerplate.javaspring.core.task.TaskRepositoryPort;

/** Activate with {@code SPRING_PROFILES_ACTIVE=postgres}. */
@Repository
@Profile("postgres")
public class PostgresTaskRepository implements TaskRepositoryPort {

    private final TaskJpaRepository jpaRepository;
    private final TaskJpaMapper mapper;

    public PostgresTaskRepository(TaskJpaRepository jpaRepository, TaskJpaMapper mapper) {
        this.jpaRepository = jpaRepository;
        this.mapper = mapper;
    }

    @Override
    public Task create(Task task) {
        return mapper.toDomain(jpaRepository.save(mapper.toEntity(task)));
    }

    @Override
    public Optional<Task> get(String id) {
        return jpaRepository.findById(id).map(mapper::toDomain);
    }

    @Override
    public List<Task> list() {
        return jpaRepository.findAll().stream().map(mapper::toDomain).toList();
    }

    @Override
    public Task update(Task task) {
        if (!jpaRepository.existsById(task.id())) {
            throw new TaskNotFoundException(task.id());
        }
        return mapper.toDomain(jpaRepository.save(mapper.toEntity(task)));
    }

    @Override
    public void delete(String id) {
        if (!jpaRepository.existsById(id)) {
            throw new TaskNotFoundException(id);
        }
        jpaRepository.deleteById(id);
    }
}
