package ai.tasktrooper.boilerplate.quarkus.platform.resilience;

import java.time.Duration;
import java.time.Instant;
import java.util.Map;
import java.util.concurrent.ConcurrentHashMap;

/**
 * Fixed-window rate limiter keyed by client. In-memory on purpose: it stops one
 * client from hammering one instance without dragging Redis into a starter
 * project. Behind more than one instance, back it with a shared store (or let
 * the gateway do the limiting) — the call sites do not change.
 */
public final class RateLimiter {

    public record Decision(boolean allowed, int remaining, long retryAfterSeconds) { }

    private record Window(int count, Instant resetAt) { }

    private final int max;
    private final Duration window;
    private final Map<String, Window> windows = new ConcurrentHashMap<>();

    public RateLimiter(int max, Duration window) {
        this.max = max > 0 ? max : 120;
        this.window = window != null && !window.isZero() ? window : Duration.ofMinutes(1);
    }

    public Decision check(String key) {
        Instant now = Instant.now();
        Window updated = windows.compute(key, (k, current) -> {
            if (current == null || !current.resetAt().isAfter(now)) {
                return new Window(1, now.plus(window));
            }
            return new Window(current.count() + 1, current.resetAt());
        });
        if (windows.size() > 10_000) {
            windows.entrySet().removeIf(e -> !e.getValue().resetAt().isAfter(now));
        }
        int remaining = Math.max(0, max - updated.count());
        long retryAfter = Math.max(1, Duration.between(now, updated.resetAt()).toSeconds());
        return new Decision(updated.count() <= max, remaining, retryAfter);
    }
}
