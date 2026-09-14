package ai.tasktrooper.boilerplate.javaspring;

import org.junit.jupiter.api.Test;
import org.springframework.boot.test.context.SpringBootTest;

/**
 * Guards against exactly the failure this repo hit once: with
 * postgres/mongo/cassandra starters all on one classpath, Spring Boot
 * tries to autoconfigure all three simultaneously and refuses to start
 * unless the unused ones are excluded (see application.yml and
 * .ai/architecture.md#config). No DB profile active here — this is the
 * in-memory default.
 */
@SpringBootTest(properties = "grpc.server.port=-1")
class SmokeTest {
    @Test
    void contextLoadsWithNoDbProfileActive() {
    }
}
