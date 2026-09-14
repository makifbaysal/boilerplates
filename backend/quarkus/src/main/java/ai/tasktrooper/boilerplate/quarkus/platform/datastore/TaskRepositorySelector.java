package ai.tasktrooper.boilerplate.quarkus.platform.datastore;

import jakarta.enterprise.context.ApplicationScoped;
import jakarta.enterprise.inject.Instance;
import jakarta.enterprise.inject.Produces;

import org.eclipse.microprofile.config.inject.ConfigProperty;

import ai.tasktrooper.boilerplate.quarkus.adapter.outbound.cassandra.CassandraTaskRepository;
import ai.tasktrooper.boilerplate.quarkus.adapter.outbound.memory.InMemoryTaskRepository;
import ai.tasktrooper.boilerplate.quarkus.adapter.outbound.mongo.MongoTaskRepository;
import ai.tasktrooper.boilerplate.quarkus.adapter.outbound.postgres.PostgresTaskRepository;
import ai.tasktrooper.boilerplate.quarkus.core.task.TaskRepositoryPort;

/**
 * The Quarkus equivalent of java-spring's {@code @Profile} on each DB
 * adapter — picks exactly one TaskRepositoryPort implementation at
 * runtime, based on {@code app.datastore}. {@code Instance<X>}
 * injection points are lazy: only the branch actually {@code .get()}'d
 * is instantiated, so e.g. CassandraTaskRepository (and its CqlSession,
 * which really connects) never gets constructed unless
 * {@code app.datastore=cassandra}.
 *
 * <p>Each concrete repository is {@code @Typed} to just itself (see
 * their javadoc) so this producer is the *only* bean exposing
 * TaskRepositoryPort — without that, this producer plus four concrete
 * beans all implementing the same interface would itself be ambiguous.
 */
@ApplicationScoped
public class TaskRepositorySelector {

    @Produces
    @ApplicationScoped
    public TaskRepositoryPort select(
            @ConfigProperty(name = "app.datastore", defaultValue = "memory") String datastore,
            Instance<InMemoryTaskRepository> memory,
            Instance<PostgresTaskRepository> postgres,
            Instance<MongoTaskRepository> mongo,
            Instance<CassandraTaskRepository> cassandra) {
        return switch (datastore) {
            case "postgres" -> postgres.get();
            case "mongo" -> mongo.get();
            case "cassandra" -> cassandra.get();
            default -> memory.get();
        };
    }
}
