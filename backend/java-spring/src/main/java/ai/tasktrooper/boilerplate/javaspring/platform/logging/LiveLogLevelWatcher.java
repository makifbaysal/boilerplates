package ai.tasktrooper.boilerplate.javaspring.platform.logging;

import java.io.IOException;
import java.io.InputStream;
import java.nio.file.FileSystems;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.StandardWatchEventKinds;
import java.nio.file.WatchEvent;
import java.nio.file.WatchKey;
import java.nio.file.WatchService;
import java.util.Map;
import java.util.Optional;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.boot.context.event.ApplicationReadyEvent;
import org.springframework.context.ApplicationListener;
import org.springframework.core.env.Environment;
import org.springframework.stereotype.Component;
import org.yaml.snakeyaml.Yaml;

import ch.qos.logback.classic.Level;

/**
 * Dynamic config reload — the Java equivalent of backend/go-fiber's
 * config.Loader.Watch wired into a slog.LevelVar. Watches the active
 * environment profile's application-&lt;profile&gt;.yml on the real
 * filesystem (CONFIG_DIR if set — see Dockerfile — otherwise
 * src/main/resources for local `mvn spring-boot:run`) and applies a
 * changed {@code logging.level.root} without a restart.
 *
 * <p>This only watches a real file on disk. Once packaged, the same
 * bytes are also embedded in the jar's classpath (how Spring resolved
 * them at startup) — this component intentionally targets the external
 * copy, matching go-fiber's "config lives in a directory next to the
 * binary" design rather than trying to rewrite an embedded jar entry.
 */
@Component
public class LiveLogLevelWatcher implements ApplicationListener<ApplicationReadyEvent> {

    private static final Logger log = LoggerFactory.getLogger(LiveLogLevelWatcher.class);

    private final Environment environment;

    public LiveLogLevelWatcher(Environment environment) {
        this.environment = environment;
    }

    @Override
    public void onApplicationEvent(ApplicationReadyEvent event) {
        Optional<Path> watched = resolveEnvProfileFile();
        if (watched.isEmpty()) {
            log.info("no dev/staging/prod profile file found on disk, dynamic log-level reload disabled");
            return;
        }
        Path file = watched.get();
        Thread watcherThread = new Thread(() -> watch(file), "live-log-level-watcher");
        watcherThread.setDaemon(true);
        watcherThread.start();
        log.info("watching {} for dynamic log-level reload", file);
    }

    private Optional<Path> resolveEnvProfileFile() {
        String configDir = System.getenv().getOrDefault("CONFIG_DIR", "src/main/resources");
        for (String profile : environment.getActiveProfiles()) {
            if (!profile.equals("dev") && !profile.equals("staging") && !profile.equals("prod")) {
                continue;
            }
            Path candidate = Path.of(configDir, "application-" + profile + ".yml");
            if (Files.isRegularFile(candidate)) {
                return Optional.of(candidate);
            }
        }
        return Optional.empty();
    }

    private void watch(Path file) {
        Path dir = file.toAbsolutePath().getParent();
        try (WatchService watchService = FileSystems.getDefault().newWatchService()) {
            dir.register(watchService, StandardWatchEventKinds.ENTRY_MODIFY);
            while (true) {
                WatchKey key = watchService.take();
                for (WatchEvent<?> watchEvent : key.pollEvents()) {
                    Path changed = dir.resolve((Path) watchEvent.context());
                    if (changed.getFileName().equals(file.getFileName())) {
                        reload(file);
                    }
                }
                if (!key.reset()) {
                    break;
                }
            }
        } catch (IOException | InterruptedException e) {
            log.error("log-level watch stopped", e);
            Thread.currentThread().interrupt();
        }
    }

    // package-private so LiveLogLevelWatcherTest can drive a reload
    // directly without spinning up a real WatchService + Spring context.
    void reload(Path file) {
        try (InputStream in = Files.newInputStream(file)) {
            Map<String, Object> yaml = new Yaml().load(in);
            String level = digLevel(yaml);
            if (level == null) {
                return;
            }
            var rootLogger = (ch.qos.logback.classic.Logger) LoggerFactory.getLogger(Logger.ROOT_LOGGER_NAME);
            rootLogger.setLevel(Level.valueOf(level));
            log.info("log level reloaded from {}: root={}", file, level);
        } catch (IOException e) {
            log.error("failed to reload {}, keeping previous log level", file, e);
        }
    }

    @SuppressWarnings("unchecked")
    private static String digLevel(Map<String, Object> yaml) {
        Object logging = yaml == null ? null : yaml.get("logging");
        if (!(logging instanceof Map<?, ?> loggingMap)) {
            return null;
        }
        Object level = ((Map<String, Object>) loggingMap).get("level");
        if (!(level instanceof Map<?, ?> levelMap)) {
            return null;
        }
        Object root = ((Map<String, Object>) levelMap).get("root");
        return root == null ? null : root.toString();
    }
}
