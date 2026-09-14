package ai.tasktrooper.boilerplate.quarkus.adapter.outbound.cassandra;

import java.util.ArrayList;
import java.util.List;
import java.util.Optional;

import jakarta.enterprise.context.ApplicationScoped;
import jakarta.enterprise.inject.Typed;
import jakarta.inject.Inject;

import com.datastax.oss.driver.api.core.CqlSession;
import com.datastax.oss.driver.api.core.cql.Row;

import ai.tasktrooper.boilerplate.quarkus.core.task.Task;
import ai.tasktrooper.boilerplate.quarkus.core.task.TaskNotFoundException;
import ai.tasktrooper.boilerplate.quarkus.core.task.TaskRepositoryPort;

/** Activate with {@code app.datastore=cassandra}. See schema.cql for
 * the table this expects, and TaskRepositorySelector /
 * InMemoryTaskRepository's javadoc for why {@code @Typed} is here. */
@ApplicationScoped
@Typed(CassandraTaskRepository.class)
public class CassandraTaskRepository implements TaskRepositoryPort {

    private final CqlSession session;

    @Inject
    public CassandraTaskRepository(CqlSession session) {
        this.session = session;
    }

    private static Task toDomain(Row row) {
        return new Task(
                row.getString("id"),
                row.getString("title"),
                row.getBoolean("done"),
                row.getInstant("created_at"),
                row.getInstant("updated_at"));
    }

    @Override
    public Task create(Task task) {
        session.execute(
                "INSERT INTO tasks (id, title, done, created_at, updated_at) VALUES (?, ?, ?, ?, ?)",
                task.id(), task.title(), task.done(), task.createdAt(), task.updatedAt());
        return task;
    }

    @Override
    public Optional<Task> get(String id) {
        Row row = session.execute("SELECT id, title, done, created_at, updated_at FROM tasks WHERE id = ?", id).one();
        return row == null ? Optional.empty() : Optional.of(toDomain(row));
    }

    @Override
    public List<Task> list() {
        List<Task> out = new ArrayList<>();
        session.execute("SELECT id, title, done, created_at, updated_at FROM tasks")
                .forEach(row -> out.add(toDomain(row)));
        return out;
    }

    // Cassandra's UPDATE is an upsert at the CQL level — check existence
    // first so this port's "update must already exist" contract holds
    // the same way it does for postgres/mongo.
    @Override
    public Task update(Task task) {
        if (get(task.id()).isEmpty()) {
            throw new TaskNotFoundException(task.id());
        }
        session.execute(
                "UPDATE tasks SET title = ?, done = ?, updated_at = ? WHERE id = ?",
                task.title(), task.done(), task.updatedAt(), task.id());
        return task;
    }

    @Override
    public void delete(String id) {
        if (get(id).isEmpty()) {
            throw new TaskNotFoundException(id);
        }
        session.execute("DELETE FROM tasks WHERE id = ?", id);
    }
}
