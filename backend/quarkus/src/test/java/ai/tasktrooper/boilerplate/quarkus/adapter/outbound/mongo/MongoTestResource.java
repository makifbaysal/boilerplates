package ai.tasktrooper.boilerplate.quarkus.adapter.outbound.mongo;

import java.util.Map;

import org.testcontainers.containers.MongoDBContainer;

import io.quarkus.test.common.QuarkusTestResourceLifecycleManager;

public class MongoTestResource implements QuarkusTestResourceLifecycleManager {

    private MongoDBContainer container;

    @Override
    public Map<String, String> start() {
        container = new MongoDBContainer("mongo:8");
        container.start();
        return Map.of(
                "app.datastore", "mongo",
                "quarkus.mongodb.connection-string", container.getConnectionString(),
                "quarkus.mongodb.database", "quarkus_it");
    }

    @Override
    public void stop() {
        if (container != null) {
            container.stop();
        }
    }
}
