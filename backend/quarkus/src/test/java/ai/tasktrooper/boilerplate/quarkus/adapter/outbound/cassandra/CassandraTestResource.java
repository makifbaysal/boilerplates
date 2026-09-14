package ai.tasktrooper.boilerplate.quarkus.adapter.outbound.cassandra;

import java.util.Map;

import org.testcontainers.cassandra.CassandraContainer;

import com.datastax.oss.driver.api.core.CqlSession;

import io.quarkus.test.common.QuarkusTestResourceLifecycleManager;

public class CassandraTestResource implements QuarkusTestResourceLifecycleManager {

    private CassandraContainer container;

    @Override
    public Map<String, String> start() {
        container = new CassandraContainer("cassandra:5");
        container.start();

        try (CqlSession session = CqlSession.builder()
                .addContactPoint(container.getContactPoint())
                .withLocalDatacenter(container.getLocalDatacenter())
                .build()) {
            session.execute("CREATE KEYSPACE IF NOT EXISTS quarkus_it WITH replication = "
                    + "{'class':'SimpleStrategy','replication_factor':1}");
            session.execute("CREATE TABLE IF NOT EXISTS quarkus_it.tasks ("
                    + "id text PRIMARY KEY, title text, done boolean, created_at timestamp, updated_at timestamp)");
        }

        return Map.of(
                "app.datastore", "cassandra",
                "app.cassandra.contact-point", container.getContactPoint().getHostString(),
                "app.cassandra.port", String.valueOf(container.getContactPoint().getPort()),
                "app.cassandra.local-datacenter", container.getLocalDatacenter(),
                "app.cassandra.keyspace", "quarkus_it");
    }

    @Override
    public void stop() {
        if (container != null) {
            container.stop();
        }
    }
}
