package ai.tasktrooper.boilerplate.javaspring.platform.logging;

import static org.assertj.core.api.Assertions.assertThat;

import java.io.IOException;
import java.nio.file.Files;
import java.nio.file.Path;

import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.io.TempDir;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.mock.env.MockEnvironment;

import ch.qos.logback.classic.Level;

/**
 * Drives LiveLogLevelWatcher#reload directly against a real file on
 * disk — proves the parse-and-apply half of the mechanism without the
 * timing flakiness of waiting on a real WatchService (that plumbing is
 * ~15 lines of java.nio.file, not worth a slow/flaky test on its own).
 */
class LiveLogLevelWatcherTest {

    @Test
    void reloadAppliesNewRootLevelFromFile(@TempDir Path dir) throws IOException {
        Path file = dir.resolve("application-dev.yml");
        Files.writeString(file, "logging:\n  level:\n    root: warn\n");

        var watcher = new LiveLogLevelWatcher(new MockEnvironment().withProperty("dummy", "x"));
        var rootLogger = (ch.qos.logback.classic.Logger) LoggerFactory.getLogger(Logger.ROOT_LOGGER_NAME);
        Level before = rootLogger.getLevel();
        try {
            watcher.reload(file);
            assertThat(rootLogger.getLevel()).isEqualTo(Level.WARN);

            Files.writeString(file, "logging:\n  level:\n    root: debug\n");
            watcher.reload(file);
            assertThat(rootLogger.getLevel()).isEqualTo(Level.DEBUG);
        } finally {
            rootLogger.setLevel(before);
        }
    }
}
