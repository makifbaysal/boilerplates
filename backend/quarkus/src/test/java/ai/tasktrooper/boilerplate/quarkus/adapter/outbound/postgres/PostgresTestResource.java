package ai.tasktrooper.boilerplate.quarkus.adapter.outbound.postgres;

import java.util.Map;

import org.testcontainers.containers.PostgreSQLContainer;

import io.quarkus.test.common.QuarkusTestResourceLifecycleManager;

public class PostgresTestResource implements QuarkusTestResourceLifecycleManager {

    private PostgreSQLContainer<?> container;

    @Override
    public Map<String, String> start() {
        container = new PostgreSQLContainer<>("postgres:17-alpine");
        container.start();
        return Map.of(
                "app.datastore", "postgres",
                "quarkus.hibernate-orm.active", "true",
                "quarkus.hibernate-orm.database.generation", "drop-and-create",
                "quarkus.datasource.db-kind", "postgresql",
                "quarkus.datasource.jdbc.url", container.getJdbcUrl(),
                "quarkus.datasource.username", container.getUsername(),
                "quarkus.datasource.password", container.getPassword());
    }

    @Override
    public void stop() {
        if (container != null) {
            container.stop();
        }
    }
}
