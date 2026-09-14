package ai.tasktrooper.boilerplate.javaspring.adapter.outbound.mongo;

import static org.assertj.core.api.Assertions.assertThat;

import java.time.Instant;

import org.junit.jupiter.api.Test;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.test.context.SpringBootTest;
import org.springframework.test.context.ActiveProfiles;
import org.springframework.test.context.DynamicPropertyRegistry;
import org.springframework.test.context.DynamicPropertySource;
import org.testcontainers.containers.MongoDBContainer;
import org.testcontainers.junit.jupiter.Container;
import org.testcontainers.junit.jupiter.Testcontainers;

import ai.tasktrooper.boilerplate.javaspring.core.task.Task;
import ai.tasktrooper.boilerplate.javaspring.core.task.TaskRepositoryPort;

/** Run with {@code mvn verify} — see PostgresTaskRepositoryIT's javadoc. */
@Testcontainers
@SpringBootTest(properties = "grpc.server.port=-1")
@ActiveProfiles("mongo")
class MongoTaskRepositoryIT {

    @Container
    static MongoDBContainer mongo = new MongoDBContainer("mongo:8");

    @DynamicPropertySource
    static void registerProperties(DynamicPropertyRegistry registry) {
        registry.add("spring.data.mongodb.uri", () -> mongo.getConnectionString() + "/worker_it");
    }

    @Autowired
    private TaskRepositoryPort repository;

    @Test
    void createGetUpdateDelete() {
        Instant now = Instant.now();
        Task created = repository.create(new Task("it-1", "buy milk", false, now, now));

        assertThat(repository.get("it-1")).isPresent();
        assertThat(repository.list()).isNotEmpty();

        Task updated = repository.update(created.withDone(true, Instant.now()));
        assertThat(updated.done()).isTrue();

        repository.delete("it-1");
        assertThat(repository.get("it-1")).isEmpty();
    }
}
