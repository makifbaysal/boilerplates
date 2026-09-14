package ai.tasktrooper.boilerplate.javaspring.adapter.outbound.cassandra;

import static org.assertj.core.api.Assertions.assertThat;

import java.time.Instant;

import org.junit.jupiter.api.Test;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.test.context.SpringBootTest;
import org.springframework.test.context.ActiveProfiles;
import org.springframework.test.context.DynamicPropertyRegistry;
import org.springframework.test.context.DynamicPropertySource;
import org.testcontainers.cassandra.CassandraContainer;
import org.testcontainers.junit.jupiter.Container;
import org.testcontainers.junit.jupiter.Testcontainers;

import com.datastax.oss.driver.api.core.CqlSession;

import ai.tasktrooper.boilerplate.javaspring.core.task.Task;
import ai.tasktrooper.boilerplate.javaspring.core.task.TaskRepositoryPort;

/** Run with {@code mvn verify} — see PostgresTaskRepositoryIT's javadoc. */
@Testcontainers
@SpringBootTest(properties = "grpc.server.port=-1")
@ActiveProfiles("cassandra")
class CassandraTaskRepositoryIT {

    @Container
    static CassandraContainer cassandra = new CassandraContainer("cassandra:5");

    @DynamicPropertySource
    static void registerProperties(DynamicPropertyRegistry registry) {
        registry.add("spring.cassandra.contact-points", () -> cassandra.getContactPoint().getHostString());
        registry.add("spring.cassandra.port", () -> cassandra.getContactPoint().getPort());
        registry.add("spring.cassandra.local-datacenter", cassandra::getLocalDatacenter);
        registry.add("spring.cassandra.keyspace-name", () -> "worker_it");
        registry.add("spring.cassandra.schema-action", () -> "create-if-not-exists");
    }

    @org.junit.jupiter.api.BeforeAll
    static void createKeyspace() {
        try (CqlSession session = CqlSession.builder()
                .addContactPoint(cassandra.getContactPoint())
                .withLocalDatacenter(cassandra.getLocalDatacenter())
                .build()) {
            session.execute("CREATE KEYSPACE IF NOT EXISTS worker_it WITH replication = "
                    + "{'class':'SimpleStrategy','replication_factor':1}");
        }
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
