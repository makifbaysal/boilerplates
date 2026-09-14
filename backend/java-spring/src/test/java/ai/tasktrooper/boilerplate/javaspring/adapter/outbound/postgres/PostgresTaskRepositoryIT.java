package ai.tasktrooper.boilerplate.javaspring.adapter.outbound.postgres;

import static org.assertj.core.api.Assertions.assertThat;
import static org.assertj.core.api.Assertions.assertThatThrownBy;

import java.time.Instant;

import org.junit.jupiter.api.Test;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.test.context.SpringBootTest;
import org.springframework.test.context.ActiveProfiles;
import org.springframework.test.context.DynamicPropertyRegistry;
import org.springframework.test.context.DynamicPropertySource;
import org.testcontainers.containers.PostgreSQLContainer;
import org.testcontainers.junit.jupiter.Container;
import org.testcontainers.junit.jupiter.Testcontainers;

import ai.tasktrooper.boilerplate.javaspring.core.task.Task;
import ai.tasktrooper.boilerplate.javaspring.core.task.TaskNotFoundException;
import ai.tasktrooper.boilerplate.javaspring.core.task.TaskRepositoryPort;

/**
 * Run with {@code mvn verify} (failsafe picks up *IT.java, plain
 * {@code mvn test} does not — the Maven equivalent of go-fiber's
 * {@code //go:build integration} skip pattern; see .ai/testing.md).
 * Needs a real Docker daemon.
 */
@Testcontainers
@SpringBootTest(properties = "grpc.server.port=-1") // this IT only cares about the repository, skip binding a real gRPC port
@ActiveProfiles("postgres")
class PostgresTaskRepositoryIT {

    @Container
    static PostgreSQLContainer<?> postgres = new PostgreSQLContainer<>("postgres:17-alpine");

    @DynamicPropertySource
    static void registerProperties(DynamicPropertyRegistry registry) {
        registry.add("spring.datasource.url", postgres::getJdbcUrl);
        registry.add("spring.datasource.username", postgres::getUsername);
        registry.add("spring.datasource.password", postgres::getPassword);
        registry.add("spring.jpa.hibernate.ddl-auto", () -> "create-drop");
    }

    @Autowired
    private TaskRepositoryPort repository;

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
        assertThatThrownBy(() -> repository.delete("it-1")).isInstanceOf(TaskNotFoundException.class);
    }
}
