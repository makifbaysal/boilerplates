package ai.tasktrooper.boilerplate.quarkus.adapter.outbound.postgres;

import static org.assertj.core.api.Assertions.assertThat;

import java.time.Instant;

import jakarta.inject.Inject;

import org.junit.jupiter.api.Test;

import io.quarkus.test.common.QuarkusTestResource;
import io.quarkus.test.junit.QuarkusTest;

import ai.tasktrooper.boilerplate.quarkus.core.task.Task;
import ai.tasktrooper.boilerplate.quarkus.core.task.TaskRepositoryPort;

/**
 * Run with {@code mvn verify} — see .ai/testing.md#integration-tests.
 * PostgresTestResource starts a real Testcontainers Postgres and flips
 * app.datastore=postgres + quarkus.hibernate-orm.active=true for just
 * this test class. Needs a real Docker daemon.
 */
@QuarkusTest
@QuarkusTestResource(value = PostgresTestResource.class, restrictToAnnotatedClass = true)
class PostgresTaskRepositoryIT {

    @Inject
    TaskRepositoryPort repository;

    @Test
    void createGetUpdateDelete() {
        Instant now = Instant.now();
        Task created = repository.create(new Task("it-1", "buy milk", false, now, now));
        assertThat(created.id()).isEqualTo("it-1");

        assertThat(repository.get("it-1")).isPresent();
        assertThat(repository.list()).isNotEmpty();

        Task updated = repository.update(created.withDone(true, Instant.now()));
        assertThat(updated.done()).isTrue();

        repository.delete("it-1");
        assertThat(repository.get("it-1")).isEmpty();
    }
}
