package ai.tasktrooper.boilerplate.quarkus.adapter.outbound.cassandra;

import java.net.InetSocketAddress;

import jakarta.enterprise.context.ApplicationScoped;
import jakarta.enterprise.inject.Disposes;
import jakarta.enterprise.inject.Produces;

import org.eclipse.microprofile.config.inject.ConfigProperty;

import com.datastax.oss.driver.api.core.CqlSession;

/**
 * No first-party/well-maintained Quarkus Cassandra extension exists
 * (DataStax's own Quarkus integration is stale, pinned to a much older
 * Quarkus baseline) — wired in directly as a plain CDI producer instead,
 * same driver go-fiber/worker/java-spring use. Only actually connects
 * when CassandraTaskRepository is instantiated, which only happens when
 * TaskRepositorySelector picks it — see that class.
 */
@ApplicationScoped
public class CqlSessionProducer {

    @Produces
    @ApplicationScoped
    public CqlSession produceCqlSession(
            @ConfigProperty(name = "app.cassandra.contact-point", defaultValue = "localhost") String contactPoint,
            @ConfigProperty(name = "app.cassandra.port", defaultValue = "9042") int port,
            @ConfigProperty(name = "app.cassandra.local-datacenter", defaultValue = "datacenter1") String localDatacenter,
            @ConfigProperty(name = "app.cassandra.keyspace", defaultValue = "quarkus") String keyspace) {
        return CqlSession.builder()
                .addContactPoint(new InetSocketAddress(contactPoint, port))
                .withLocalDatacenter(localDatacenter)
                .withKeyspace(keyspace)
                .build();
    }

    public void closeCqlSession(@Disposes CqlSession session) {
        session.close();
    }
}
