package ai.tasktrooper.boilerplate.javaspring.platform.resilience;

import java.time.Duration;
import java.time.Instant;
import java.util.concurrent.Callable;

/**
 * Dependency-free circuit breaker for outbound calls (database, HTTP API,
 * broker). Use one instance per dependency: a shared breaker makes a healthy
 * dependency unavailable because a different one is broken.
 *
 * <p>Kept hand-written rather than pulling in Resilience4j/SmallRye so the
 * boilerplate does not impose a resilience library on the project that copies
 * it — swap the internals for one of those without touching call sites.
 */
public final class CircuitBreaker {

    public enum State { CLOSED, OPEN, HALF_OPEN }

    /** Thrown instead of calling the dependency while the circuit is open. */
    public static final class OpenCircuitException extends RuntimeException {
        public OpenCircuitException(String name) {
            super("circuit breaker is open: " + name);
        }
    }

    private final String name;
    private final int failureThreshold;
    private final Duration openDuration;
    private final int halfOpenMaxCalls;

    private State state = State.CLOSED;
    private int failures;
    private int halfOpenCalls;
    private Instant openedAt = Instant.EPOCH;

    public CircuitBreaker(String name) {
        this(name, 5, Duration.ofSeconds(30), 1);
    }

    public CircuitBreaker(String name, int failureThreshold, Duration openDuration, int halfOpenMaxCalls) {
        this.name = name;
        this.failureThreshold = failureThreshold > 0 ? failureThreshold : 5;
        this.openDuration = openDuration != null ? openDuration : Duration.ofSeconds(30);
        this.halfOpenMaxCalls = halfOpenMaxCalls > 0 ? halfOpenMaxCalls : 1;
    }

    public synchronized State state() {
        refresh();
        return state;
    }

    /** Runs the call unless the circuit is open. */
    public <T> T call(Callable<T> action) throws Exception {
        acquire();
        try {
            T result = action.call();
            onSuccess();
            return result;
        } catch (Exception e) {
            onFailure();
            throw e;
        }
    }

    public void run(Runnable action) {
        try {
            call(() -> {
                action.run();
                return null;
            });
        } catch (RuntimeException e) {
            throw e;
        } catch (Exception e) {
            throw new IllegalStateException(e);
        }
    }

    private synchronized void acquire() {
        refresh();
        if (state == State.OPEN) {
            throw new OpenCircuitException(name);
        }
        if (state == State.HALF_OPEN) {
            if (halfOpenCalls >= halfOpenMaxCalls) {
                throw new OpenCircuitException(name);
            }
            halfOpenCalls++;
        }
    }

    private synchronized void refresh() {
        if (state == State.OPEN && Duration.between(openedAt, Instant.now()).compareTo(openDuration) >= 0) {
            state = State.HALF_OPEN;
            halfOpenCalls = 0;
        }
    }

    private synchronized void onSuccess() {
        if (state == State.HALF_OPEN) {
            state = State.CLOSED;
            halfOpenCalls = 0;
        }
        failures = 0;
    }

    private synchronized void onFailure() {
        if (state == State.HALF_OPEN) {
            state = State.OPEN;
            openedAt = Instant.now();
            halfOpenCalls = 0;
            return;
        }
        failures++;
        if (failures >= failureThreshold) {
            state = State.OPEN;
            openedAt = Instant.now();
        }
    }
}
